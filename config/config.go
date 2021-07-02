
package config

import (
	"fmt"
	"log"
	"path"

	"github.com/adrg/xdg"
	"github.com/pelletier/go-toml"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var (
	defaultConfigDir  string
	DefaultConfigPath string
	defaultStorageDir string
	defaultConfig     = map[string]interface{}{
		"VarFromFile": "ViperTest From file",
	}
)

func init() {
	var err error
	DefaultConfigPath, err = xdg.ConfigFile("apitesting/apitesting.toml")
	if err != nil {
		log.Fatal(err)
	}
	defaultConfigDir = path.Dir(DefaultConfigPath)
	defaultStorageDir = path.Join(xdg.DataHome, "apitesting")

	err = createDefaultFile(afero.NewOsFs())
	if err != nil {
		log.Fatal(err)
	}
}

type initConfigErr struct {
	s string
}

func (e *initConfigErr) Error() string {
	return e.s
}

func newInitConfigErr(err error) error {
	return &initConfigErr{
		s: fmt.Sprintf("Init config error: %s", err.Error()),
	}
}

func createDefaultFile(fs afero.Fs) error {
	err := fs.MkdirAll(defaultConfigDir, 0755)
	if err != nil {
		return err
	}
	fs.MkdirAll(defaultStorageDir, 0755)
	if err != nil {
		return err
	}

	exist, err := afero.Exists(fs, DefaultConfigPath)
	if err != nil {
		return err
	}

	if !exist {
		handle, err := fs.Create(DefaultConfigPath)
		if err != nil {
			return err
		}
		defer handle.Close()
		t, err := toml.TreeFromMap(defaultConfig)
		if err != nil {
			return err
		}
		handle.WriteString(t.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func InitConfig(fs afero.Fs, configPath string) error {
	if configPath == "" {
		viper.SetConfigFile(DefaultConfigPath)
	} else {
		exist, err := afero.Exists(fs, configPath)
		if err != nil {
			return newInitConfigErr(err)
		}
		if !exist {
			return &initConfigErr{
				s: fmt.Sprintf("Init config error: %s not exist", configPath),
			}
		}
		viper.SetConfigFile(configPath)
	}
	return viper.ReadInConfig()
}
