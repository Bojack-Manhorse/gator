package config

import (
	"os"

	"github.com/Bojack-Manhorse/pokedexcli/gator/utils"
)

func getPath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := homePath + "/.gatorconfig.json"

	return path, nil
}

func Read() (Config, error) {
	path, err := getPath()
	if err != nil {
		return Config{}, err
	}
	jsonFile, err := os.Open(path)

	if err != nil {
		return Config{}, err
	}

	defer jsonFile.Close()

	return utils.DecodeReader[Config](jsonFile)
}
