package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	selfsigned "github.com/wolfeidau/golang-self-signed-tls"
)

type RequestData struct {
	// "confirm", "password", "text"
	Type string `json:"type"`
	// Message for dialog window
	Message string `json:"message"`
	// Device name used to automatically search for a password in a folder named 'code'
	// and retrieve the password if a match is found.
	Device string `json:"device"`
	// For example, the client may request a password for a specific application,
	// such as "LUKS" (e.g., when a password is needed after a device reboot and disk reattachment).
	// Or, for example, request confirmation for a USB flash drive inserted into the device;
	// you can specify "usb" as an additional hint.
	Code string `json:"code"`
}

func makeNewServer(server string, certIps []string, callback func(w http.ResponseWriter, r *http.Request)) {

	tlsCert, err := generateSelfSignedCert(certIps)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr: server,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
		},
		Handler:      ipFilterMiddleware(http.HandlerFunc(callback)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Starting server on " + server)
	log.Println("Trusted IPs: " + strings.Join(allowedIPs, ", "))
	err = srv.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}
}

func generateSelfSignedCert(addresses []string) (tls.Certificate, error) {
	result, err := selfsigned.GenerateCert(
		selfsigned.Hosts(addresses),
		selfsigned.RSABits(2048),
		selfsigned.ValidFor(365*24*time.Hour),
	)
	if err != nil {
		return tls.Certificate{}, err
	}

	cert, err := tls.X509KeyPair(result.PublicCert, result.PrivateKey)
	if err != nil {
		return tls.Certificate{}, err
	}

	log.Println("Certificate fingerprint:", result.Fingerprint)
	return cert, nil
}

func ipFilterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Invalid remote address", http.StatusForbidden)
			return
		}
		ip = net.ParseIP(ip).String()
		if !slices.Contains(allowedIPs, ip) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isValidPort(port string) bool {

	p, err := strconv.Atoi(port)
	if err != nil {
		return false
	}

	return p > 0 && p <= 65535
}
