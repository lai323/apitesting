package entity

import (
	"apitesting/config"
	"testing"
)

func TestInstall(t *testing.T) {
	err := config.InitConfig("../config.toml")
	if err != nil {
		panic(err)
	}
	err = Install()
	if err != nil {
		panic(err)
	}
}
