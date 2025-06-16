package main

import (
	"path/filepath"
	"time"

	"github.com/gen2brain/dlgs"
)

type guiRequest struct {
	data   RequestData
	result chan string
}

var guiRequests = make(chan guiRequest)

func guiWorker() {
	for req := range guiRequests {
		var response string

		done := make(chan struct{})
		go func() {
			switch req.data.Type {
			case "confirm":
				ok, err := dlgs.Question(req.data.Device, req.data.Message, true)
				if err != nil || !ok {
					response = ""
				} else if ok {
					response = "1"
				}
			case "password":
				password, _ := findPassword(config.Passwords, req.data.Code, req.data.Device)
				if len(password) > 0 {
					ok, err := dlgs.Question(req.data.Device, filepath.Join(config.Passwords, req.data.Code, req.data.Device)+" → 🔑 → "+req.data.Code, true)
					if err != nil || !ok {
						response = ""
					} else if ok {
						response = password
					}
				} else {
					password, ok, err := dlgs.Password(req.data.Device, req.data.Message)
					if err != nil || !ok {
						response = ""
					} else {
						response = password
						go func() {
							ok, err := dlgs.Question(req.data.Device, "🔑 → "+filepath.Join(config.Passwords, req.data.Code, req.data.Device), true)
							if err != nil || !ok {
								response = ""
							} else if ok {
								setPassword(config.Passwords, req.data.Code, req.data.Device, password)
							}
						}()
					}
				}
			case "text":
				input, ok, err := dlgs.Entry(req.data.Device, req.data.Message, "")
				if err != nil || !ok {
					response = ""
				} else {
					response = input
				}
			default:
				response = ""
			}
			close(done)
		}()

		select {
		case <-done:
			req.result <- response
		case <-time.After(30 * time.Second):
			req.result <- ""
		}
	}
}
