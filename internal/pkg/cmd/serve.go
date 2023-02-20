package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zltl/nydus-auth/internal/pkg/api"
	"github.com/zltl/nydus-auth/internal/pkg/db"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the nydus-auth server",
	Long: `Start the nydus-auth server, which will listen on the port specified by the
environment variable NYDUS_AUTH_PORT. If NYDUS_AUTH_PORT is not set, it will
listen on port 8080.`,
	Run: func(cmd *cobra.Command, args []string) {
		execServe()
	},
}

func execServe() {
	logrus.Println("start serve")
	initConfig()
	initLogger()
	db.Ctx.Open("")

	srv := api.State{}
	srv.Start()

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("$PWD")
		viper.AddConfigPath("/etc/nydus-auth/")
		viper.AddConfigPath("$HOME/.nydus-auth/")
		viper.SetConfigName("conf")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		DisableQuote:  true,
	})
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)
}
