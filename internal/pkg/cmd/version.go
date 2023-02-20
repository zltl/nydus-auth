package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zltl/nydus-auth/pkg/version"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of nydus-auth",
	Long:  `All software has versions. This is nydus-auth's`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(fmt.Sprintf("nydus-auth version %s", version.Version))
	},
}
