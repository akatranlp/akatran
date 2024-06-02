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

	dnsRepo "github.com/akatranlp/akatran/internal/dns"
	"github.com/akatranlp/akatran/internal/spinner"
	"github.com/spf13/cobra"
)

var jsonOutput bool
var listRecordType string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [flags] domain",
	Short: "List all DNS records of selected domain",
	Args:  cobra.ExactArgs(1),
	Long: `With the subcommands you can list all records of a given domain.
You have to provide the domain as an argument. And can choose between table and json output. 
For example:

  akatran dns [--token <cloudflare-token>] [--provider <cloudflare>] list example.com 

  -------------------------------------------
  | TYPE  |       NAME      |    CONTENT    |
  -------------------------------------------
  | A     |     example.com | 198.51.100.10 |
  | AAAA  |     example.com |   2001:DB8::1 |
  | CNAME | www.example.com |   example.com |
  -------------------------------------------
	

  akatran dns list example.com --json

  [{"name":"example.com","type":"A","content":"192.51.100.10"},
  {"name":"example.com","type":"AAAA","content":"2001:DB8::1"},
  {"name":"www.example.com","type":"CNAME","content":"example.com"}]
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SetErrPrefix("Error: [DNS - LIST] - ")

		domain := args[0]

		repo, err := dnsRepo.GetRepoFromViperOrFlag(domain, provider, token)
		if err != nil {
			return err
		}

		spinner.Start()
		defer spinner.Stop()

		cmd.Println("Listing DNS records for", domain)

		var dnsRecords dnsRepo.DnsRecordList
		switch listRecordType {
		case "":
			dnsRecords, err = repo.ListRecords(cmd.Context())
			if err != nil {
				return err
			}
		case "A":
			fallthrough
		case "AAAA":
			fallthrough
		case "CNAME":
			dnsRecords, err = repo.ListRecords(cmd.Context(), listRecordType)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid record type")
		}

		spinner.Stop()

		if jsonOutput {
			cmd.Println(dnsRecords.AsJsonString())
		} else {
			cmd.Println(dnsRecords.AsTableString())
		}
		return nil
	},
}

func init() {
	DnsCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&listRecordType, "type", "t", "", "The type of the DNS record")
	listCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output as JSON")
}
