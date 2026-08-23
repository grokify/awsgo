package main

import (
	"fmt"
	"log"

	"github.com/grokify/awsgo/rdsutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	File string `short:"f" long:"file" description:"RDS DescribeEngineDefaultParameters/DescribeDBParameters JSON export to read" required:"true"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	params, err := rdsutil.ParametersResponseReadFile(opts.File)
	logutil.FatalErr(err)
	fmtutil.MustPrintJSON(params)
	fmt.Printf("PARAMETER COUNT (%d)\n", len(params.Parameters))
}
