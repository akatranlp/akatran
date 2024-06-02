/*
Copyright © 2024 Fabian Petersen <fabian@nf-petersen.de>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package dns

import (
	"fmt"
	"net/url"
	"strings"

	dnsRepo "github.com/akatranlp/akatran/internal/dns"
	"github.com/akatranlp/akatran/internal/spinner"
	"github.com/akatranlp/akatran/internal/utils"
	"github.com/spf13/cobra"
)

var recordType string
var recordContent string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [flags] dns_record",
	Short: "Create a new DNS record",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SetErrPrefix("Error: [DNS - CREATE] - ")

		dnsRecord := args[0]

		parts := strings.Split(dnsRecord, ".")
		if len(parts) < 2 {
			return fmt.Errorf("invalid dns record")
		}

		domain := strings.Join(parts[len(parts)-2:], ".")

		repo, err := dnsRepo.GetRepoFromViperOrFlag(domain, provider, token)
		if err != nil {
			return err
		}

		switch recordType {
		case "A":
			ip := utils.GetIPv4Address(recordContent)
			if ip == nil {
				return fmt.Errorf("invalid IPv4 address")
			}
			recordContent = ip.String()

		case "AAAA":
			ip := utils.GetIPv6Address(recordContent)
			if ip == nil {
				return fmt.Errorf("invalid IPv6 address")
			}
			recordContent = ip.String()
		case "CNAME":
			if recordContent == "" {
				return fmt.Errorf("content is required for CNAME records")
			}
			domain, err := url.Parse(recordContent)
			if err != nil {
				return err
			}
			recordContent = domain.Hostname()
		default:
			return fmt.Errorf("invalid record type")
		}

		spinner.Start()
		defer spinner.Stop()

		if err := repo.CreateRecord(cmd.Context(), dnsRepo.DnsRecord{
			Name:    dnsRecord,
			Type:    recordType,
			Content: recordContent,
		}); err != nil {
			return err
		}

		spinner.Stop()
		cmd.Printf("DNS record %s created!\n", dnsRecord)
		return nil
	},
}

func init() {
	DnsCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&recordType, "type", "t", "A", "The type of the DNS record")
	createCmd.Flags().StringVarP(&recordContent, "content", "c", "", "The content of the DNS record")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
