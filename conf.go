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

	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(exePath)

	return filepath.Join(dir, "config.json"), nil

}

// func configInit() (Config, error) {

// 	configFile, err := getConfigFilePath()
// 	if err != nil {
// 		return Config{}, err
// 	}

// 	if _, err := os.Stat(configFile); os.IsNotExist(err) {
// 		var devices string
// 		var passwords string
// 		fmt.Print("Configuration file not found. Please enter the path to the folder: ")
// 		_, err := fmt.Scanln(&devices)
// 		if err != nil {
// 			fmt.Println("Error reading input:", err)
// 			return Config{}, err
// 		}

// 		if err := createConfig(configFile, devices, passwords); err != nil {
// 			fmt.Println("Error creating configuration file:", err)
// 			return Config{}, err
// 		}
// 		fmt.Println("Created configuration file:", configFile)
// 	}

// 	config, err := readConfig(configFile)
// 	if err != nil {
// 		fmt.Println("Error reading configuration file:", err)
// 		return Config{}, err
// 	}

// 	return config, nil

// }

func configInit() (Config, error) {

	configFile, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Print("Config not found. Enter path to save config file (or press Enter to save in this dir): ")
		var userPath string
		fmt.Scanln(&userPath)
		if userPath != "" {
			configFile = userPath
		}
		var trusted, server, passwords string
		fmt.Print("Enter trusted IPs path:\n" +
			"- can be a file with allowed IPs\n" +
			"- or a folder (or multiple folders) for recursive search (e.g. /path/**/filename)\n> ")
		fmt.Scanln(&trusted)
		fmt.Print("Enter server (ip:port) address: ")
		fmt.Scanln(&server)
		fmt.Print("Enter folder to store passwords: ")
		fmt.Scanln(&passwords)
		if err := createConfig(configFile, trusted, server, passwords); err != nil {
			return Config{}, err
		}
		fmt.Println("Created config at:", configFile)
	}

	return readConfig(configFile)

}

func readConfig(configFile string) (Config, error) {

	file, err := os.Open(configFile)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil

}

func createConfig(configFile, trusted, server, passwords string) error {

	file, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(Config{Trusted: trusted, Server: server, Passwords: passwords})

}
