package iamutil

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/grokify/awsgo/config"
)

const (
	ActionSTSAssumeRole = "sts:AssumeRole"
)

var (
	ErrIAMClientNotSet  = errors.New("aws iam client not set")
	ErrPolicyNameNotSet = errors.New("aws policy name not set")
	ErrRoleNameNotSet   = errors.New("aws role name not set")
	ErrUserNameNotSet   = errors.New("aws user name not set")
)

type IAMService struct {
	AWSIAMClient *iam.Client
	Policies     *PolicyService
	Roles        *RoleService
	Users        *UserService
}

func NewIAMService(ctx context.Context, cfg *config.AWSConfig, optFns ...func(*iam.Options)) (*IAMService, error) {
	awsV2Cfg, err := cfg.ConfigV2(ctx)
	if err != nil {
		return nil, err
	}
	iamClient := iam.NewFromConfig(awsV2Cfg, optFns...)
	return &IAMService{
		AWSIAMClient: iamClient,
		Policies:     &PolicyService{AWSIAMClient: iamClient},
		Roles:        &RoleService{AWSIAMClient: iamClient},
		Users:        &UserService{AWSIAMClient: iamClient},
	}, nil
}
