package state

import (
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// GetHln returns the .hln dir.
func GetHln() string {
	var (
		home string
		err  error
	)

	homeDir := viper.GetString("home")
	if homeDir == "" {
		homeDir, err = homedir.Dir()
		if err != nil {
			panic(err)
		}
	}
	home = path.Join(homeDir, ".hln")
	return home
}
