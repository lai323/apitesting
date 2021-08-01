package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func InitConfig(configPath string) error {
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("InitConfig error: config file (%s) not exist", configPath)
		}
		return fmt.Errorf("InitConfig error: %s", err.Error())
	}

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("InitConfig error: %s", err.Error())
	}
	return nil
}
