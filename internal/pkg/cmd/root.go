package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zltl/nydus-auth/pkg/version"
)

var rootCmd = &cobra.Command{
	Use:     "nydus-auth",
	Short:   "nydus-auth is a OAuth2 server",
	Long:    `nydus-auth is a OAuth2 server, which implements the following OAuth2 flows.`,
	Version: version.Version,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var cfgFile string

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile,
		"config",
		"c",
		"",
		"config file (default is ./conf.ini)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
