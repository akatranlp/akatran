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
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/akatranlp/akatran/cmd/dns"
	"github.com/akatranlp/akatran/internal/viper"
	"github.com/akatranlp/akatran/pkg/bytesize"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var cfgFile string
var storageSize bytesize.ByteSize = 100 * bytesize.GB
var ram bytesize.ByteSize = 2 * bytesize.GiB

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "akatran",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", viper.Get("size"))

		var cfg Config
		err := viper.UnmarshalExact(&cfg)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(cfg.StorageSize, uint64(cfg.StorageSize))
		fmt.Println(cfg.StorageSize.Format("%.03f", "GB"))

		fmt.Println(cfg.StorageSize.FromMB())

		fmt.Println(cfg.Ram, uint64(cfg.Ram))
		fmt.Println(cfg.Ram.Format("%.03f", "GB"))

		fmt.Println(cfg.Ram.FromMiB())

		viper.GetUint64("size")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func ExecuteContext(ctx context.Context) {
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func SetAppVersion(version string) {
	rootCmd.Version = version
}

func addSubCommands() {
	rootCmd.AddCommand(dns.DnsCmd)
}

func init() {
	addSubCommands()

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is $XDG_CONFIG_HOME/%s/config.yaml)", rootCmd.Name()))

	rootCmd.PersistentFlags().VarP(&storageSize, "size", "s", "storage size")
	viper.BindPFlag("size", rootCmd.PersistentFlags().Lookup("size"))
	rootCmd.PersistentFlags().VarP(&ram, "ram", "r", "ram size")
	viper.BindPFlag("ram", rootCmd.PersistentFlags().Lookup("ram"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		configDir, err := os.UserConfigDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(path.Join(configDir, rootCmd.Name()))
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	if err := godotenv.Load(); err == nil {
		fmt.Fprintln(os.Stderr, "Using .env file")
	}

	viper.SetEnvPrefix("AKATRAN")
	viper.SetEnvKeyReplacer(strings.NewReplacer("::", "_"))

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

type Config struct {
	StorageSize bytesize.ByteSize `mapstructure:"size"`
	Ram         bytesize.ByteSize `mapstructure:"ram"`
}
