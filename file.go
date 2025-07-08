package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func findTrustedIpAddress(pattern string) ([]string, error) {

	files, err := findTrustedFiles(parseTrustedPattern(pattern))
	if err != nil {
		return nil, err
	}
	ipSet := make(map[string]struct{})
	for _, f := range files {
		lines, err := readTrustedFile(f)
		if err != nil {
			return nil, err
		}
		for _, ip := range lines {
			ipSet[ip] = struct{}{}
		}
	}
	var ips []string
	for ip := range ipSet {
		ips = append(ips, ip)
	}

	return ips, nil

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

func parseTrustedPattern(pattern string) (base, name string) {

	name = filepath.Base(pattern)
	base = strings.Split(pattern, "**")[0]

	return strings.TrimRight(base, "/"), name

}

func findTrustedFiles(base, name string) ([]string, error) {

	var files []string
	err := filepath.WalkDir(base, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == name && !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err

}

func readTrustedFile(path string) ([]string, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var result []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.FieldsFunc(line, func(r rune) bool {
			return r == ' ' || r == ',' || r == '\t' || r == '\n'
		})
		for _, f := range fields {
			if f != "" {
				result = append(result, f)
			}
		}
	}

	return result, scanner.Err()

}
