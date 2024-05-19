package iamutil

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/grokify/mogo/pointer"
	"github.com/micahhausler/aws-iam-policy/policy"
)

type RoleService struct {
	AWSIAMClient *iam.Client
}

// Create used used to create a role.
// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/iam
// https://docs.aws.amazon.com/code-library/latest/ug/go_2_iam_code_examples.html
func (svc RoleService) Create(ctx context.Context, params CreateRoleInput, optFns ...func(*iam.Options)) (*types.Role, *iam.CreateRoleOutput, error) {
	if svc.AWSIAMClient == nil {
		return nil, nil, ErrIAMClientNotSet
	} else if input, err := params.Request(); err != nil {
		return nil, nil, err
	} else if result, err := svc.AWSIAMClient.CreateRole(ctx, input, optFns...); err != nil {
		return nil, nil, err
	} else {
		return result.Role, result, nil
	}
}

func CreatePolicyDocumentCreateRoleWithTrustedUserARN(roleName string, trustedUserARN string) policy.Policy {
	trustPolicy := policy.Policy{
		Version: policy.VersionLatest,
		Statements: policy.NewStatementOrSlice([]policy.Statement{
			{
				Effect:    policy.EffectAllow,
				Principal: policy.NewAWSPrincipal(trustedUserARN),
				Action:    policy.NewStringOrSlice(true, ActionSTSAssumeRole),
			},
		}...),
	}
	return trustPolicy
}

/*

trustPolicy := PolicyDocument{
	Version:   "2012-10-17",
	Statement: []PolicyStatement{{
		Effect: "Allow",
		Principal: map[string]string{"AWS": trustedUserArn},
		Action: []string{"sts:AssumeRole"},
	}},
}
*/

// CreateRoleInput is a representation of the AWS v2 SDK Create Role Input where
// the AssumeRolePolicyDocument required as a `policy.Policy` struct.
type CreateRoleInput struct {
	RoleName                 string        // required
	AssumeRolePolicyDocument policy.Policy // required
	iam.CreateRoleInput
}

func (input CreateRoleInput) Request() (*iam.CreateRoleInput, error) {
	if strings.TrimSpace(input.RoleName) == "" {
		return nil, ErrRoleNameNotSet
	} else if b, err := json.Marshal(input.AssumeRolePolicyDocument); err != nil {
		return nil, err
	} else {
		return &iam.CreateRoleInput{
			RoleName:                 pointer.Pointer(input.RoleName),
			AssumeRolePolicyDocument: pointer.Pointer(string(b)),
			Description:              input.CreateRoleInput.Description,
			MaxSessionDuration:       input.CreateRoleInput.MaxSessionDuration,
			Path:                     input.CreateRoleInput.Path,
			PermissionsBoundary:      input.CreateRoleInput.PermissionsBoundary,
			Tags:                     input.CreateRoleInput.Tags}, nil
	}
}
