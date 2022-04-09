package state

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

// GetTemp returns the path to the temp directory.
func GetTemp() string {
	var temp string
	tempDir := viper.GetString("temp_home")
	if tempDir == "" {
		tempDir = os.TempDir()
	}
	temp = path.Join(tempDir, "heighliner")
	return temp
}

// CleanTemp removes all temporary resources.
func CleanTemp() {
	if err := os.RemoveAll(GetTemp()); err != nil {
		panic(err)
	}
}
