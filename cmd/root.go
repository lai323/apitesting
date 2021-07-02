
package cmd

import (
	"fmt"
	"os"
	"apitesting/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configPath string
	rootCmd    = &cobra.Command{
		Use:   "readygo module_name",
		Short: "create empty project with cobra and spf13",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			// Do Stuff Here
		},
	}
	subCmd = &cobra.Command{
		Use:   "subcmd",
		Short: "sub command example",
		Run: func(cmd *cobra.Command, args []string) {
			// cmd.Help()
			fmt.Println(viper.GetString("var"))
			fmt.Println(viper.GetString("VarFromFile"))
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
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", fmt.Sprintf("config file (default is %s)", config.DefaultConfigPath))
	rootCmd.PersistentFlags().String("var", "ViperTest var", "use Viper for configuration")
	viper.BindPFlag("var", rootCmd.PersistentFlags().Lookup("var"))
	rootCmd.AddCommand(subCmd)
}

func initConfig() {
	err := config.InitConfig(afero.NewOsFs(), configPath)
	if err != nil {
		panic(err)
	}
}

