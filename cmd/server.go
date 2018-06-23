package cmd

import (
	"fmt"
	"os"

	"github.com/nickwu241/simply-do/server"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the simply-do server",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := server.NewServer()
		if err != nil {
			fmt.Printf("error initializing server: %v\n", err)
			os.Exit(2)
		}
		server.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
