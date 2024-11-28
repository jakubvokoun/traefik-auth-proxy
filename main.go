package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := flag.String("port", "8080", "Port to run the HTTP server on")
	flag.Parse()

	allowedIPsEnv := os.Getenv("ALLOWED_IPS")
	if allowedIPsEnv == "" {
		fmt.Println("Environment variable ALLOWED_IPS is not set. Exiting.")
		os.Exit(1)
	}
	allowedIPs := strings.Split(allowedIPsEnv, ",")

	http.HandleFunc("/traefik", func(w http.ResponseWriter, r *http.Request) {
		xForwardedFor := r.Header.Get("X-Forwarded-For")

		if xForwardedFor == "" {
			http.Error(w, "X-Forwarded-For header not found", http.StatusBadRequest)
			return
		}
		fmt.Printf("X-Forwarded-For IP: %s\n", xForwardedFor)

		ip := net.ParseIP(xForwardedFor)
		if ip == nil {
			http.Error(w, "Invalid IP address in X-Forwarded-For header", http.StatusBadRequest)
			return
		}

		matched := false
		for _, allowed := range allowedIPs {
			_, network, err := net.ParseCIDR(allowed)
			if err == nil {
				// CIDR block
				if network.Contains(ip) {
					matched = true
					break
				}
			} else {
				// Individual IP address
				if allowed == ip.String() {
					matched = true
					break
				}
			}
		}

		if matched {
			fmt.Fprintf(w, "X-Forwarded-For: %s matches allowed IPs or CIDR blocks\n", xForwardedFor)
		} else {
			http.Error(w, "Forbidden: IP  not in allowed IPs or CIDR blocks", http.StatusForbidden)
		}
	})

	fmt.Printf("Server is running on port %s...\n", *port)
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
