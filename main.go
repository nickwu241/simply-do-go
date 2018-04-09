package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/nickwu241/simply-do/server"
)

func main() {
	server, err := server.NewServer()
	if err != nil {
		fmt.Printf("error initializing server: %v\n", err)
		os.Exit(2)
	}
	// Use PORT from environment variables if it's set. Needed for Heroku.
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port, err := strconv.Atoi(portEnv)
		if err != nil {
			fmt.Printf("error converting PORT environment to integer: %v\n", err)
			os.Exit(2)
		}
		server.Run(fmt.Sprintf(":%d", port))
	} else {
		server.Run(":8080")
	}
}
