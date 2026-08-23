package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "awsgo",
	Short: "AWS utility CLI backed by the awsgo Go libraries",
	Long: `awsgo is a command-line toolbox for AWS operational tasks, built on
top of the awsgo Go libraries (rds, and more over time).`,
}
