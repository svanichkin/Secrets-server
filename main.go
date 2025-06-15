package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
)

var allowedIPs []string

func main() {

	var err error
	config, err = configInit()
	if err != nil {
		fmt.Println("Failed to initialize config:", err)
		return
	}
	allowedIPs, err = findTrustedIpAddress(config.Trusted)
	if err != nil {
		fmt.Println("Please add a file named 'trusted' with trusted IPs to the folder '" + config.Trusted + "' then 'trusted' files are searched recursively in this folder and its subfolders.")
		return
	}
	server, err := readServer(config.Server)
	if err != nil {
		fmt.Println("Please create the file '" + config.Trusted + "' and add IP:port to start the server.")
		return
	}
	host, port, err := net.SplitHostPort(server)
	if err != nil || net.ParseIP(host) == nil || !isValidPort(port) {
		fmt.Println("Please check the file '" + config.Trusted + "' and check IP:port.")
		return
	}

	go guiWorker()
	makeNewServer(server, []string{"localhost", "127.0.0.1", host}, dialogHandler)

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
