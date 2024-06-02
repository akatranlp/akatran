package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/akatranlp/akatran/internal/utils"
)

type DnsRecord struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type DnsRecordList []DnsRecord

func (d DnsRecordList) AsTableString() string {
	paddingName := 4
	paddingContent := 7
	for _, record := range d {
		paddingName = max(len(record.Name), paddingName)
		paddingContent = max(len(record.Content), paddingContent)
	}

	var builder strings.Builder

	spacer := strings.Repeat("-", 10+5+paddingName+paddingContent)
	fmt.Fprintln(&builder, spacer)

	typeString := "TYPE "

	newPaddingName := paddingName - 4
	nameString := fmt.Sprintf("%sNAME%s", strings.Repeat(" ", int(math.Ceil(float64(newPaddingName)/2))), strings.Repeat(" ", int(math.Floor(float64(newPaddingName)/2))))

	newPaddingContent := paddingContent - 7
	contentString := fmt.Sprintf("%sCONTENT%s", strings.Repeat(" ", int(math.Ceil(float64(newPaddingContent)/2))), strings.Repeat(" ", int(math.Floor(float64(newPaddingContent)/2))))

	fmt.Fprintf(&builder, "| %s | %s | %s |\n", typeString, nameString, contentString)
	fmt.Fprintln(&builder, spacer)
	for _, record := range d {
		fmt.Fprintf(&builder, "| %-5s | %*s | %*s |\n", record.Type, paddingName, record.Name, paddingContent, record.Content)
	}
	fmt.Fprintln(&builder, spacer)

	return builder.String()
}

func (d DnsRecordList) AsJsonString() string {
	var builder strings.Builder
	json.NewEncoder(&builder).Encode(d)

	return builder.String()
}

type DnsRepository interface {
	ListRecords(ctx context.Context, types ...string) (DnsRecordList, error)
	CreateRecord(ctx context.Context, record DnsRecord) error
	UpdateRecord(ctx context.Context, record DnsRecord) error
	DeleteRecord(ctx context.Context, record DnsRecord) (DnsRecordList, error)
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
		case a.Type == "A":
			return -1
		case b.Type == "A":
			return 1
		case a.Type == "AAAA":
			return -1
		default:
			return 1
		}
	})

	return utils.Map(sortedRecords, func(record SortDnsRecords) DnsRecord {
		return record.DnsRecord
	})
}
