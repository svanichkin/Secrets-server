package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var allowedIPs []string

func main() {

	var err error
	config, err = configInit()
	if err != nil {
		fmt.Println("Failed to initialize config:", err)
		return
	}
	allowedIPs, err := findTrustedIpAddress(config.Trusted)
	if err != nil {
		fmt.Printf("Error searching for trusted IPs files in '%s': %v\n", config.Trusted, err)
		return
	}
	if len(allowedIPs) == 0 {
		fmt.Printf("No trusted IPs found. Please add IP in '%s' file.\n", config.Trusted)
		return
	}
	server := config.Server
	host, port, err := net.SplitHostPort(server)
	if err != nil || net.ParseIP(host) == nil || !isValidPort(port) {
		fmt.Println("Please check config path '" + config.Trusted + "' and verify IP:port entries.")
		return
	}

	go guiWorker()
	makeNewServer(server, allowedIPs, dialogHandler)

	// If binary updated - restart

	exePath, _ := os.Executable()
	info, _ := os.Stat(exePath)
	lastModTime := info.ModTime()
	for {
		time.Sleep(30 * time.Second)
		info, err := os.Stat(exePath)
		if err != nil {
			continue
		}
		if info.ModTime() != lastModTime {
			fmt.Println("Binary updated, stop service...")
			exec.Command(exePath).Start()
			os.Exit(0)
		}
	}

}

func dialogHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var reqData RequestData
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	resultChan := make(chan string)
	guiRequests <- guiRequest{data: reqData, result: resultChan}
	answer := <-resultChan
	if answer == "" {
		switch reqData.Type {
		case "confirm":
			http.Error(w, "User denied confirmation", http.StatusBadRequest)
		case "password", "text":
			http.Error(w, "User cancelled input", http.StatusBadRequest)
		default:
			http.Error(w, "Unsupported request type", http.StatusBadRequest)
		}
		return
	}
	log.Println("Answer received from GUI:", answer)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(answer))

}
