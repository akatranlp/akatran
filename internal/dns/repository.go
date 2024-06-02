package dns

import (
	"context"
	"slices"
	"strings"

	"github.com/akatranlp/akatran/internal/utils"
)

type dnsTypes string

const (
	A     dnsTypes = "A"
	AAAA  dnsTypes = "AAAA"
	CNAME dnsTypes = "CNAME"
)

type DnsRecord struct {
	Name    string   `json:"name"`
	Type    dnsTypes `json:"type"`
	Content string   `json:"content"`
}

type DnsRepository interface {
	ListRecords(ctx context.Context) ([]DnsRecord, error)
	CreateRecord(ctx context.Context, record DnsRecord) error
}

func sortDnsRecords(records []DnsRecord) []DnsRecord {
	type SortDnsRecords struct {
		DnsRecord
		nameParts []string
	}
	var sortedRecords = utils.Map(records, func(record DnsRecord) SortDnsRecords {
		nameparts := strings.Split(record.Name, ".")
		slices.Reverse(nameparts)

		return SortDnsRecords{
			DnsRecord: record,
			nameParts: nameparts,
		}
	})

	slices.SortFunc(sortedRecords, func(a, b SortDnsRecords) int {
		switch {
		case a.Type == b.Type:
			return slices.Compare(a.nameParts, b.nameParts)
		case a.Type == A:
			return -1
		case b.Type == A:
			return 1
		case a.Type == AAAA:
			return -1
		default:
			return 1
		}
	})

	return utils.Map(sortedRecords, func(record SortDnsRecords) DnsRecord {
		return record.DnsRecord
	})
}
