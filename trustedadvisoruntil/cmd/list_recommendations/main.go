package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grokify/awsgo/config"
	trustedadvisorutil "github.com/grokify/awsgo/trustedadvisoruntil"
	"github.com/grokify/mogo/log/logutil"
	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	CredsFile string `short:"c" long:"credentialsfile" description:"Create subscription" required:"true"`
	CredsKey  string `short:"k" long:"key" description:"Create subscription" required:"true"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.AWSConfigCredentialsSetFile(opts.CredsFile, opts.CredsKey)
	logutil.FatalErr(err)

	tcsvc, err := trustedadvisorutil.NewTrustedAdvisorService(context.Background(), cfg)
	logutil.FatalErr(err)

	recs, err := tcsvc.ListRecommendations(context.Background(), nil)
	logutil.FatalErr(err)

	fmt.Printf("REC COUNT (%d)\n", len(recs.RecommendationSummaries))

	fmt.Println("DONE")
}
