package main

import "github.com/spf13/cobra"

var rdsCmd = &cobra.Command{
	Use:   "rds",
	Short: "RDS utilities",
}

func init() {
	rootCmd.AddCommand(rdsCmd)
}
