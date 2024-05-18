package iamutil

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/grokify/awsgo/config"
)

type IAMService struct {
	AWSIAMClient  *iam.Client
	PolicyService *PolicyService
}

func NewIAMService(ctx context.Context, cfg *config.AWSConfig, optFns ...func(*iam.Options)) (*IAMService, error) {
	awsV2Cfg, err := cfg.ConfigV2(ctx)
	if err != nil {
		return nil, err
	}
	iamClient := iam.NewFromConfig(awsV2Cfg, optFns...)
	return &IAMService{
		AWSIAMClient:  iamClient,
		PolicyService: &PolicyService{AWSIAMClient: iamClient}}, nil
}
