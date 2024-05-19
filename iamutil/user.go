package iamutil

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type UserService struct {
	AWSIAMClient *iam.Client
}

// Create used used to create a user.
// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/iam
// https://docs.aws.amazon.com/code-library/latest/ug/go_2_iam_code_examples.html
func (svc UserService) Create(ctx context.Context, params CreateUserInput, optFns ...func(*iam.Options)) (*types.User, *iam.CreateUserOutput, error) {
	if svc.AWSIAMClient == nil {
		return nil, nil, ErrIAMClientNotSet
	} else if input, err := params.Request(); err != nil {
		return nil, nil, err
	} else if result, err := svc.AWSIAMClient.CreateUser(ctx, input, optFns...); err != nil {
		return nil, nil, err
	} else {
		return result.User, result, nil
	}
}

func (svc UserService) CreateSimple(ctx context.Context, username string, params *CreateUserInput, optFns ...func(*iam.Options)) (*types.User, *iam.CreateUserOutput, error) {
	if params == nil {
		params = &CreateUserInput{}
	}
	if username != "" {
		params.UserName = username
	}
	if svc.AWSIAMClient == nil {
		return nil, nil, ErrIAMClientNotSet
	} else if input, err := params.Request(); err != nil {
		return nil, nil, err
	} else if result, err := svc.AWSIAMClient.CreateUser(ctx, input, optFns...); err != nil {
		return nil, nil, err
	} else {
		return result.User, result, nil
	}
}

// CreateRoleInput is a representation of the AWS v2 SDK Create Uuser Input where
// the UserName property is required.
type CreateUserInput struct {
	UserName string // required
	iam.CreateUserInput
}

func (input CreateUserInput) Request() (*iam.CreateUserInput, error) {
	if strings.TrimSpace(input.UserName) == "" {
		return nil, ErrUserNameNotSet
	} else if b, err := json.Marshal(input); err != nil {
		return nil, err
	} else {
		var out *iam.CreateUserInput
		return out, json.Unmarshal(b, out)
	}
}
