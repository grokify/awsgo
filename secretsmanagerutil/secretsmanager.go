package secretsmanagerutil

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/grokify/mogo/pointer"
)

/*
func NewCache() (*secretcache.Cache, error) {
	// "github.com/aws/aws-secretsmanager-caching-go/secretcache"
	return secretcache.New()
}
*/

type Client struct {
	awsSecretsManagerClient *secretsmanager.Client
}

func (c *Client) GetSecretStringSimple(secretID string) (string, error) {
	return c.GetSecretString(context.Background(), secretID, "", "")
}

func (c *Client) GetSecretString(ctx context.Context, secretID, versionID, versionStage string, optFns ...func(*secretsmanager.Options)) (string, error) {
	opts := secretsmanager.GetSecretValueInput{}
	if secretID != "" {
		opts.SecretId = pointer.Pointer(secretID)
	}
	if versionID != "" {
		opts.VersionId = pointer.Pointer(versionID)
	}
	if versionStage != "" {
		opts.VersionStage = pointer.Pointer(versionStage)
	}
	res, err := c.awsSecretsManagerClient.GetSecretValue(ctx, &opts)
	if err != nil {
		return "", err
	} else if res.SecretString == nil {
		return "", errors.New("secret string is nil")
	}
	return pointer.Dereference(res.SecretString), nil
}

func (c *Client) UpdateSecretStringSimple(secretID, secretValue string) (versionID string, err error) {
	opts := secretsmanager.UpdateSecretInput{}
	if secretID != "" {
		opts.SecretId = pointer.Pointer(secretID)
	}
	if secretValue != "" {
		opts.SecretString = pointer.Pointer(secretValue)
	}
	res, err := c.awsSecretsManagerClient.UpdateSecret(context.Background(), &opts)
	if err != nil {
		return "", err
	} else if res.VersionId == nil {
		return "", errors.New("versionID not set")
	}
	return pointer.Dereference(res.VersionId), nil
}
