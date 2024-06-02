package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const CloudflareProvider string = "cloudflare"

var (
	ErrGetZoneIDFailes   = fmt.Errorf("failed to get zone id from domain")
	ErrDomainNotFound    = fmt.Errorf("domain not found")
	ErrListZonesFailed   = fmt.Errorf("failed to list zones")
	ErrListRecordsFailed = fmt.Errorf("failed to list records")
)

type CloudflareRepo struct {
	domain string
	token  string
	client *http.Client
}

func NewCloudflareRepo(domain, token string, client *http.Client) *CloudflareRepo {
	return &CloudflareRepo{
		domain: domain,
		token:  token,
		client: client,
	}
}

type listRecordsResponse struct {
	Result []DnsRecord `json:"result"`
}

type DnsZone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *CloudflareRepo) getZoneIDFromDomain(ctx context.Context, domain string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.cloudflare.com/client/v4/zones", nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("name", domain)
	q.Add("status", "active")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", ErrListZonesFailed
	}

	var zones struct {
		Result []DnsZone `json:"result"`
	}

	if err := json.NewDecoder(res.Body).Decode(&zones); err != nil {
		return "", err
	}

	if len(zones.Result) == 0 {
		return "", ErrDomainNotFound
	}

	result := zones.Result[0]
	if result.Name != domain {
		return "", ErrDomainNotFound
	}

	return result.ID, nil
}

func (c *CloudflareRepo) ListRecords(ctx context.Context) ([]DnsRecord, error) {
	zoneID, err := c.getZoneIDFromDomain(ctx, c.domain)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	var dnsRecords listRecordsResponse
	if err := json.NewDecoder(res.Body).Decode(&dnsRecords); err != nil {
		return nil, err
	}

	records := make([]DnsRecord, 0)
	for _, record := range dnsRecords.Result {
		if record.Type != A && record.Type != AAAA && record.Type != CNAME {
			continue
		}
		records = append(records, record)
	}

	return sortDnsRecords(records), nil
}

func (c *CloudflareRepo) CreateRecord(ctx context.Context, record DnsRecord) error {
	zoneID, err := c.getZoneIDFromDomain(ctx, c.domain)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(struct {
		DnsRecord
		TTL int `json:"ttl"`
	}{
		DnsRecord: record,
		TTL:       1,
	}); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID), &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return err
	}
	return nil
}
