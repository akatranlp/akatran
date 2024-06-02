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
	"github.com/spf13/cobra"
)

var token string
var provider string

// DnsCmd represents the Dns command
var DnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "This command allows you to interact with DNS records",
	Long: `With the subcommands you can list, create, update or delete
given DNS-records for your domains. For example:

	akatran dns [--token <cloudflare-token>] [--provider <cloudflare>] list example.com 
	akatran dns create www.example.com --type A --content 127.0.0.1
	akatran dns update www.example.com --content 192.168.0.1
	akatran dns delete www.example.com
`,
}

func init() {
	DnsCmd.PersistentFlags().StringVar(&provider, "provider", "", "DNS provider")
	DnsCmd.PersistentFlags().StringVar(&token, "token", "", "API token")
}
