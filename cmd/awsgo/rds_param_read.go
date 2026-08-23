package main

import (
	"fmt"

	"github.com/grokify/awsgo/rdsutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/spf13/cobra"
)

var rdsParamReadFile string

var rdsParamReadCmd = &cobra.Command{
	Use:   "param-read",
	Short: "Print an RDS parameters JSON export and its parameter count",
	Long: `Reads a DescribeEngineDefaultParameters/DescribeDBParameters JSON export
and prints the parsed parameters plus a count.

Example:
  awsgo rds param-read --file db-perf-group_mysql8.json`,
	Args: cobra.NoArgs,
	RunE: runRDSParamRead,
}

func init() {
	rdsCmd.AddCommand(rdsParamReadCmd)

	rdsParamReadCmd.Flags().StringVarP(&rdsParamReadFile, "file", "f", "",
		"RDS parameters JSON file to read")
	if err := rdsParamReadCmd.MarkFlagRequired("file"); err != nil {
		panic(err)
	}
}

func runRDSParamRead(_ *cobra.Command, _ []string) error {
	params, err := rdsutil.ParametersResponseReadFile(rdsParamReadFile)
	if err != nil {
		return err
	}
	if err := fmtutil.PrintJSON(params); err != nil {
		return err
	}
	fmt.Printf("PARAMETER COUNT (%d)\n", len(params.Parameters))
	return nil
}
