package main

import (
	"bufio"
	"os"
	"path/filepath"
)

func findTrustedIpAddress(rootPath string) ([]string, error) {
	var trustedLines []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			trustedFile := filepath.Join(path, "trusted")
			if _, err := os.Stat(trustedFile); err == nil {
				file, err := os.Open(trustedFile)
				if err != nil {
					return err
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					trustedLines = append(trustedLines, scanner.Text())
				}
				if err := scanner.Err(); err != nil {
					return err
				}
			}
		}
		return nil
	})

	return trustedLines, err
}

func readServer(filePath string) (string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}

func findPassword(basePath, code, device string) (string, error) {

	filePath := filepath.Join(basePath, code, device)

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil

}

func setPassword(basePath, code, device, password string) error {

	filePath := filepath.Join(basePath, code, device)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(password)
	if err != nil {
		return err
	}

	return nil
}
