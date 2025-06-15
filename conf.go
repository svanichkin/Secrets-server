package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var config Config

type Config struct {
	Trusted   string `json:"trusted"`
	Server    string `json:"server"`
	Passwords string `json:"passwords"`
}

func getConfigFilePath() (string, error) {

	configDir := "/etc/secrets"
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil

}

func configInit() (Config, error) {

	configFile, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		var devices string
		var passwords string
		fmt.Print("Configuration file not found. Please enter the path to the folder: ")
		_, err := fmt.Scanln(&devices)
		if err != nil {
			fmt.Println("Error reading input:", err)
			return Config{}, err
		}

		if err := createConfig(configFile, devices, passwords); err != nil {
			fmt.Println("Error creating configuration file:", err)
			return Config{}, err
		}
		fmt.Println("Created configuration file:", configFile)
	}

	config, err := readConfig(configFile)
	if err != nil {
		fmt.Println("Error reading configuration file:", err)
		return Config{}, err
	}

	return config, nil

}

func readConfig(configFile string) (Config, error) {

	file, err := os.Open(configFile)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil

}

func createConfig(configFile, devices, passwords string) error {

	config := Config{Trusted: devices, Server: passwords}

	file, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(config)

}
