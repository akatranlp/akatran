package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
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
	Result []cloudflareDnsRecord `json:"result"`
}

type cloudflareDnsRecord struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
	Name    string `json:"name,omitempty"`
	Content string `json:"content,omitempty"`
	TTL     int    `json:"ttl,omitempty"`
}

type cloudflareDnsZone struct {
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
		Result []cloudflareDnsZone `json:"result"`
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

func (c *CloudflareRepo) listRecords(ctx context.Context, name string, type_ string) (*listRecordsResponse, error) {
	zoneID, err := c.getZoneIDFromDomain(ctx, c.domain)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("name", name)
	q.Add("type", type_)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to list records: %s", string(data))
	}

	var dnsRecords listRecordsResponse
	if err := json.NewDecoder(res.Body).Decode(&dnsRecords); err != nil {
		return nil, err
	}

	return &dnsRecords, nil
}

func (c *CloudflareRepo) ListRecords(ctx context.Context, types ...string) (DnsRecordList, error) {
	var err error
	for _, t := range types {
		if t != "A" && t != "AAAA" && t != "CNAME" {
			return nil, fmt.Errorf("invalid record type: %s", t)
		}
	}

	var dnsRecords *listRecordsResponse
	if len(types) == 0 || len(types) > 1 {
		dnsRecords, err = c.listRecords(ctx, "", "")
		if err != nil {
			return nil, err
		}
	} else {
		dnsRecords, err = c.listRecords(ctx, "", types[0])
		if err != nil {
			return nil, err
		}
	}

	if len(types) == 0 {
		types = []string{"A", "AAAA", "CNAME"}
	}

	records := make([]DnsRecord, 0)
	for _, record := range dnsRecords.Result {
		if !slices.Contains(types, record.Type) {
			continue
		}
		records = append(records, DnsRecord{
			Name:    record.Name,
			Type:    record.Type,
			Content: record.Content,
		})
	}

	return sortDnsRecords(records), nil
}

func (c *CloudflareRepo) CreateRecord(ctx context.Context, record DnsRecord) error {
	zoneID, err := c.getZoneIDFromDomain(ctx, c.domain)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&cloudflareDnsRecord{
		Name:    record.Name,
		Type:    record.Type,
		Content: record.Content,
		TTL:     1,
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
		data, err := io.ReadAll(res.Body)
		if err != nil {

			return err
		}
		return fmt.Errorf("failed to list records: %s", string(data))
	}
	return nil
}

func (c *CloudflareRepo) deleteSingleRecord(ctx context.Context, zoneID string, recordID string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID), nil)
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
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to delete record: %s", string(data))
	}
	return nil
}

type deleteErr struct {
	record cloudflareDnsRecord
	err    error
}

type deleteErrList []deleteErr

func (d deleteErrList) Error() string {
	var builder strings.Builder

	for _, err := range d {
		fmt.Fprintf(&builder, "failed to delete record %s: %s\n", err.record.Name, err.err)
	}

	return builder.String()
}

func (c *CloudflareRepo) DeleteRecord(ctx context.Context, record DnsRecord) (DnsRecordList, error) {
	zoneID, err := c.getZoneIDFromDomain(ctx, c.domain)
	if err != nil {
		return nil, err
	}

	records, err := c.listRecords(ctx, record.Name, record.Type)
	if err != nil {
		return nil, err
	}

	if len(records.Result) == 0 {
		return nil, fmt.Errorf("record not found")
	}

	successFullDeleted := make([]DnsRecord, 0)
	errs := make(deleteErrList, 0)

	for _, record := range records.Result {
		if err := c.deleteSingleRecord(ctx, zoneID, record.ID); err != nil {
			errs = append(errs, deleteErr{
				record: record,
				err:    err,
			})
		}
		successFullDeleted = append(successFullDeleted, DnsRecord{
			Name:    record.Name,
			Type:    record.Type,
			Content: record.Content,
		})
	}

	if len(errs) > 0 {
		return sortDnsRecords(successFullDeleted), errs
	}

	return sortDnsRecords(successFullDeleted), nil
}

func (c *CloudflareRepo) UpdateRecord(ctx context.Context, record DnsRecord) error {
	zoneID, err := c.getZoneIDFromDomain(ctx, c.domain)
	if err != nil {
		return err
	}

	records, err := c.listRecords(ctx, record.Name, record.Type)
	if err != nil {
		return err
	}

	if len(records.Result) == 0 {
		return nil
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&cloudflareDnsRecord{
		Name:    record.Name,
		Type:    records.Result[0].Type,
		Content: record.Content,
		TTL:     1,
	}); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, records.Result[0].ID), &buf)
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
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to delete record: %s", string(data))
	}

	return nil
}
