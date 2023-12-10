package main

import (
	"fmt"
	"log"

	"github.com/grokify/awsgo/config"
	"github.com/grokify/awsgo/pricingutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	Authfile string `short:"a" long:"authfile" description:"Goauth credentials file" required:"true"`
	Authkey  string `short:"k" long:"key" description:"Goauth credentials key" required:"true"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.AWSConfigCredentialsSetFile(opts.Authfile, opts.Authkey)
	logutil.FatalErr(err)
	prc, err := pricingutil.NewPricingClient(cfg)
	logutil.FatalErr(err)

	svcs, err := prc.Services()
	logutil.FatalErr(err)

	codes, err := svcs.ServiceCodes()
	logutil.FatalErr(err)
	fmtutil.MustPrintJSON(codes)

	fmt.Printf("NUM CODES (%d)\n", len(codes))

	prods, err := prc.Products("   AmazonRDS")
	logutil.FatalErr(err)
	fmtutil.MustPrintJSON(prods)

	fmt.Println("DONE")
}
