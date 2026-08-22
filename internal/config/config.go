package config

import "github.com/Bojack-Manhorse/pokedexcli/gator/utils"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c Config) SetUser(newUser string) error {
	currentStuct, err := Read()
	if err != nil {
		return err
	}

	currentStuct.CurrentUserName = newUser

	path, err := getPath()

	if err != nil {
		return err
	}

	err = utils.WriteToFile(path, currentStuct)

	if err != nil {
		return err
	}

	return nil
}
