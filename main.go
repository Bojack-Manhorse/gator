package main

import (
	"fmt"

	"github.com/Bojack-Manhorse/pokedexcli/gator/internal/config"
)

func main() {
	currentConfig, err := config.Read()

	if err != nil {
		fmt.Println("Cannot read config file due to error:", err)
		return
	}

	err = currentConfig.SetUser("lane")

	if err != nil {
		fmt.Println("Cannot set current config to new name due to error:", err)
		return
	}

	newConfig, err := config.Read()

	if err != nil {
		fmt.Println("Cannot read new config file due to error:", err)
		return
	}

	fmt.Println(newConfig)

}
