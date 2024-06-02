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
	"math"
	"net/http"
	"strings"

	dnsRepo "github.com/akatranlp/akatran/internal/dns"
	"github.com/akatranlp/akatran/internal/viper"
	"github.com/spf13/cobra"
)

var token string
var provider string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:        "list [flags] domain",
	Short:      "List all DNS records of selected domain",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"domain"},
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.SetErrPrefix("Error: [DNS - LIST] - ")

		domain := args[0]
		keyStart := fmt.Sprintf("dns::%s", domain)

		if sub := viper.GetString(keyStart + "::provider"); sub != "" {
			provider = sub
		} else {
			cmd.PrintErrln("domain not found in config falling back to flag params")
		}

		var repo dnsRepo.DnsRepository
		switch provider {
		case dnsRepo.CloudflareProvider:
			if sub := viper.GetString(keyStart + "::token"); sub != "" {
				token = sub
			}

			if token == "" {
				cmd.PrintErrln(cmd.ErrPrefix(), "Token not provided")
				return
			}
			repo = dnsRepo.NewCloudflareRepo(domain, token, http.DefaultClient)
		default:
			cmd.PrintErrln(cmd.ErrPrefix(), "Provider not supported", provider)
			return
		}

		dnsRecords, err := repo.ListRecords(cmd.Context())
		if err != nil {
			cmd.PrintErrln(cmd.ErrPrefix(), err)
			return
		}

		paddingName := 0
		paddingContent := 0
		for _, record := range dnsRecords {
			paddingName = max(len(record.Name), paddingName)
			paddingContent = max(len(record.Content), paddingContent)
		}

		spacer := strings.Repeat("-", 10+5+paddingName+paddingContent)
		cmd.Println(spacer)

		typeString := "TYPE "

		newPaddingName := paddingName - 4
		nameString := fmt.Sprintf("%sNAME%s", strings.Repeat(" ", int(math.Ceil(float64(newPaddingName)/2))), strings.Repeat(" ", int(math.Floor(float64(newPaddingName)/2))))

		newPaddingContent := paddingContent - 7
		contentString := fmt.Sprintf("%sCONTENT%s", strings.Repeat(" ", int(math.Ceil(float64(newPaddingContent)/2))), strings.Repeat(" ", int(math.Floor(float64(newPaddingContent)/2))))

		cmd.Printf("| %s | %s | %s |\n", typeString, nameString, contentString)
		cmd.Println(spacer)
		for _, record := range dnsRecords {
			cmd.Printf("| %-5s | %*s | %*s |\n", record.Type, paddingName, record.Name, paddingContent, record.Content)
		}
		cmd.Println(spacer)
	},
}

func init() {
	DnsCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&provider, "provider", "p", "", "DNS provider")
	listCmd.Flags().StringVarP(&token, "token", "t", "", "API token")
}
