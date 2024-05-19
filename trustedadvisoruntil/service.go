// trustedadvisorutil is a helper package for; https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/trustedadvisor
package trustedadvisorutil

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/trustedadvisor"
	"github.com/grokify/awsgo/config"
)

type TrustedAdvisorService struct {
	AWSSvcClient *trustedadvisor.Client
}

func NewTrustedAdvisorService(ctx context.Context, cfg *config.AWSConfig, optFns ...func(*trustedadvisor.Options)) (*TrustedAdvisorService, error) {
	awsV2Cfg, err := cfg.ConfigV2(ctx)
	if err != nil {
		return nil, err
	}
	taClient := trustedadvisor.NewFromConfig(awsV2Cfg, optFns...)
	return &TrustedAdvisorService{
		AWSSvcClient: taClient,
	}, nil
}
