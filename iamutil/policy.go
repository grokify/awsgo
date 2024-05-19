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

type PolicyService struct {
	AWSIAMClient *iam.Client
}

// CreatePolicy
// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/iam
// https://docs.aws.amazon.com/code-library/latest/ug/go_2_iam_code_examples.html
func (svc PolicyService) Create(ctx context.Context, params CreatePolicyInput, optFns ...func(*iam.Options)) (*types.Policy, *iam.CreatePolicyOutput, error) {
	if svc.AWSIAMClient == nil {
		return nil, nil, ErrIAMClientNotSet
	} else if polInput, err := params.Request(); err != nil {
		return nil, nil, err
	} else if result, err := svc.AWSIAMClient.CreatePolicy(ctx, polInput, optFns...); err != nil {
		return nil, nil, err
	} else {
		return result.Policy, result, nil
	}
}

func ParsePolicy(b []byte) (*policy.Policy, error) {
	pol := &policy.Policy{}
	return pol, json.Unmarshal(b, pol)
}

// CreatePolicyInput is a representation of the AWS v2 SDK Create Policy Input where
// the Policy Name is required as a string and the Policy Document is required as a
// struct (vs. string). See more at:
// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/iam#CreatePolicyInput ,
// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_grammar.html .
type CreatePolicyInput struct {
	PolicyName     string        // required
	PolicyDocument policy.Policy // required
	Description    *string
	Path           *string
	Tags           []types.Tag
}

func (input CreatePolicyInput) Request() (*iam.CreatePolicyInput, error) {
	if strings.TrimSpace(input.PolicyName) == "" {
		return nil, ErrPolicyNameNotSet
	} else if b, err := json.Marshal(input.PolicyDocument); err != nil {
		return nil, err
	} else {
		return &iam.CreatePolicyInput{
			PolicyDocument: pointer.Pointer(string(b)),
			PolicyName:     pointer.Pointer(input.PolicyName),
			Description:    input.Description,
			Path:           input.Path,
			Tags:           input.Tags}, nil
	}
}
