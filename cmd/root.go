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
		Use:   "apitesting",
		Short: "start apitesting server",
		Run: func(cmd *cobra.Command, args []string) {
			server.Serve()
		},
	}

	initDBCmd = &cobra.Command{
		Use:   "initdb",
		Short: "init database",
		Run: func(cmd *cobra.Command, args []string) {
			server.InitDB()
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
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file (required)")
	rootCmd.AddCommand(initDBCmd)
}

func initConfig() {
	if configPath == "" {
		fmt.Println("specify the configuration file with --config")
		os.Exit(1)
	}
	err := config.InitConfig(configPath)
	if err != nil {
		panic(err)
	}
}
