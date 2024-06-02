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
	"strings"

	dnsRepo "github.com/akatranlp/akatran/internal/dns"
	"github.com/akatranlp/akatran/internal/spinner"
	"github.com/spf13/cobra"
)

var deleteRecordType string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [flags] dns_record",
	Short: "Delete a DNS record",
	Args:  cobra.ExactArgs(1),
	Long: `With the subcommands you can delete the given record of your domain.
For example:

  akatran dns [--token <cloudflare-token>] [--provider <cloudflare>] delete www.example.com --type A|AAAA|CNAME
  akatran dns delete www.example.com --type A
`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		spinner.Start()
		defer spinner.Stop()

		switch deleteRecordType {
		case "A":
		case "AAAA":
		case "CNAME":
		default:
			return fmt.Errorf("invalid record type")
		}

		dnsRecords, err := repo.DeleteRecord(cmd.Context(), dnsRepo.DnsRecord{
			Name: dnsRecord,
			Type: deleteRecordType,
		})
		if err != nil {
			return err
		}

		spinner.Stop()

		cmd.Println("The following DNS records were deleted!")
		cmd.Println(dnsRecords.AsTableString())
		return nil
	},
}

func init() {
	DnsCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&deleteRecordType, "type", "t", "", "Record type")
	deleteCmd.MarkFlagRequired("type")
}
