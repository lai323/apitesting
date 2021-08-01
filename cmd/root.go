package cmd

import (
	"apitesting/config"
	"apitesting/logger"
	"apitesting/server"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	configPath string
	rootCmd    = &cobra.Command{
		Use:   "readygo module_name",
		Short: "create empty project with cobra and spf13",
		Run: func(cmd *cobra.Command, args []string) {
			server.Run()
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func init() {
	cobra.OnInitialize(initConfig, logger.InitLogger)
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file")
}

func initConfig() {
	err := config.InitConfig(configPath)
	if err != nil {
		panic(err)
	}
}
