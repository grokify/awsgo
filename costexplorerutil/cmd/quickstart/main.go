package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/grokify/awsgo/config"
	"github.com/grokify/awsgo/costexplorerutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/time/timeutil"
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

	svc, err := costexplorerutil.NewService(context.Background(), cfg)
	logutil.FatalErr(err)

	resp, err := svc.GetCostAndUsageWithResources(context.Background(),
		costexplorerutil.GetCostAndUsageWithResourcesInput{
			Filter: types.Expression{
				Dimensions: &types.DimensionValues{
					Key:    types.DimensionRegion,
					Values: []string{config.RegionUSEast1, config.RegionUSWest1},
				},
			},
			Granularity: types.GranularityDaily,
			StartDate:   time.Now().Add(-13 * timeutil.Day),
			EndDate:     time.Now(),
			Metrics:     []types.Metric{types.MetricBlendedCost, types.MetricAmortizedCost, types.MetricUsageQuantity},
		})
	logutil.FatalErr(err)

	fmtutil.PrintJSON(resp)

	fmt.Printf("DIMENSION COUNT (%d)\n", len(resp.DimensionValueAttributes))

	fmt.Println("DONE")
}
