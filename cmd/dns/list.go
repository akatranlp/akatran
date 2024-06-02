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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [flags] domain",
	Short: "List all DNS records of selected domain",
	Args:  cobra.ExactArgs(1),
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		dnsRecords, err := repo.ListRecords(cmd.Context())
		if err != nil {
			return err
		}

		spinner.Stop()

		if jsonOutput {
			fmt.Println(dnsRecords.AsJsonString())
		} else {
			fmt.Println(dnsRecords.AsTableString())
		}
		return nil
	},
}

func init() {
	DnsCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output as JSON")
}
