// costexplorerutil is a helper package for: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/costexplorer
package costexplorerutil

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/grokify/awsgo/config"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/time/timeutil"
)

type CostExplorerService struct {
	AWSSvcClient *costexplorer.Client
}

func NewService(ctx context.Context, cfg *config.AWSConfig, optFns ...func(*costexplorer.Options)) (*CostExplorerService, error) {
	awsV2Cfg, err := cfg.ConfigV2(ctx)
	if err != nil {
		return nil, err
	}
	taClient := costexplorer.NewFromConfig(awsV2Cfg, optFns...)
	return &CostExplorerService{
		AWSSvcClient: taClient,
	}, nil
}

func (svc CostExplorerService) GetCostAndUsageWithResources(ctx context.Context, params GetCostAndUsageWithResourcesInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageWithResourcesOutput, error) {
	req := params.Request()
	return svc.AWSSvcClient.GetCostAndUsageWithResources(ctx, &req, optFns...)
}

type GetCostAndUsageWithResourcesInput struct {
	Filter        types.Expression  // required
	Granularity   types.Granularity // required
	StartDate     time.Time         // required
	EndDate       time.Time         // required
	GroupBy       []types.GroupDefinition
	Metrics       []types.Metric
	NextPageToken string // optional
}

func (input GetCostAndUsageWithResourcesInput) Request() costexplorer.GetCostAndUsageWithResourcesInput {
	req := costexplorer.GetCostAndUsageWithResourcesInput{
		Filter:      pointer.Pointer(input.Filter),
		Granularity: input.Granularity,
		TimePeriod: &types.DateInterval{
			Start: pointer.Pointer(input.StartDate.Format(timeutil.RFC3339FullDate)),
			End:   pointer.Pointer(input.EndDate.Format(timeutil.RFC3339FullDate)),
		},
		GroupBy: input.GroupBy,
		Metrics: metricsToStringSlice(input.Metrics),
	}
	if strings.TrimSpace(input.NextPageToken) != "" {
		req.NextPageToken = pointer.Pointer(input.NextPageToken)
	}
	return req
}

func metricsToStringSlice(metrics []types.Metric) []string {
	var s []string
	for _, m := range metrics {
		s = append(s, string(m))
	}
	return s
}
