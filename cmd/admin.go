package cmd

import (
	"fmt"
	"os"

	"github.com/nickwu241/simply-do/server/store"
	"github.com/spf13/cobra"
)

// adminCmd represents the admin command
var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Execute admin tasks. e.g. backing up, copying data",
}

// cpCmd represents the cp command
var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "Copies a source list to destionation list",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("expecting exactly 2 arguemnts <src> <dst>")
			os.Exit(2)
		}
		src := args[0]
		dst := args[1]
		if err := store.AdminCopyData(src, dst); err != nil {
			fmt.Printf("error copying %q to %q: %v\n", src, dst, err)
			os.Exit(2)
		}
		fmt.Printf("Successfully copied %q to %q\n", src, dst)
	},
}

// dbBackupCmd represents the db-backup command
var dbBackupCmd = &cobra.Command{
	Use:   "db-backup",
	Short: "Back up the Firebase database onto Firebase storage",
	Run: func(cmd *cobra.Command, args []string) {
		if err := store.AdminSnapshot(); err != nil {
			fmt.Printf("error backing up database: %v\n", err)
			os.Exit(2)
		}
		fmt.Println("Successfully backed up Firebase database")
	},
}

func init() {
	rootCmd.AddCommand(adminCmd)
	adminCmd.AddCommand(cpCmd)
	adminCmd.AddCommand(dbBackupCmd)
}
