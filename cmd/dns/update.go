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

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update [flags] dns_record",
	Short: "Update a DNS record",
	Args:  cobra.ExactArgs(1),
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		if err := repo.UpdateRecord(cmd.Context(), dnsRepo.DnsRecord{
			Name:    dnsRecord,
			Content: recordContent,
		}); err != nil {
			cmd.PrintErrln(cmd.ErrPrefix(), err)
		}

		spinner.Stop()
		cmd.Printf("DNS record %s updated!\n", dnsRecord)

		return nil
	},
}

func init() {
	DnsCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&recordContent, "content", "c", "", "The content of the DNS record")
}
