package main

import (
	"fmt"
	"os"

	"github.com/nickwu241/simply-do/server"
)

func main() {
	server, err := server.NewServer()
	if err != nil {
		fmt.Printf("error initializing server: %v\n", err)
		os.Exit(2)
	}
	server.Run()
}
