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
package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var hashLength int

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.RangeArgs(1, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		str := args[0]

		if len(str)+4 > hashLength {
			return errors.New("the string is too long")
		}
		spaces := hashLength - len(str) - 2
		leftSpaces := spaces / 2
		rightSpaces := spaces - leftSpaces

		fmt.Printf("%s\n", strings.Repeat("#", hashLength))
		fmt.Printf("#%s%s%s#\n", strings.Repeat(" ", leftSpaces), str, strings.Repeat(" ", rightSpaces))
		fmt.Printf("%s\n", strings.Repeat("#", hashLength))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)

	formatCmd.Flags().IntVarP(&hashLength, "hash-length", "l", 45, "Length of the hashes")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// formatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// formatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
