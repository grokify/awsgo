package config

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/grokify/goauth"
)

/*

https://docs.aws.amazon.com/sdk-for-go/api/aws/session/
https://docs.aws.amazon.com/sdk-for-go/api/service/s3/

*/

const (
	CredentialsTypeEnvironment = "environment" // uses AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN (optional)
	CredentialsTypeShared      = "shared"
	CredentialsTypeStatic      = "static"

	RegionUSWest1 = "us-west-1"
)

// AWSConfigMore handles credentials from https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
type AWSConfigMore struct {
	CredentialsType string
	SharedProfile   string // `.aws/credentials`
	StaticID        string
	StaticSecret    string
	StaticToken     string
	Endpoint        string
	Region          string
	PathStyleForce  bool
}

func AWSConfigMoreCredentialBasic(creds goauth.Credentials) (*AWSConfigMore, error) {
	if creds.Type == goauth.TypeBasic {
		return &AWSConfigMore{
			CredentialsType: CredentialsTypeStatic,
			StaticID:        creds.Basic.Username,
			StaticSecret:    creds.Basic.Password,
		}, nil
	}
	return nil, errors.New("creds type not supported")
}

func (cm AWSConfigMore) Config() *aws.Config {
	cfg := &aws.Config{}
	credsType := strings.ToLower(strings.TrimSpace(cm.CredentialsType))

	if credsType == CredentialsTypeStatic {
		cfg.Credentials = credentials.NewStaticCredentials(cm.StaticID, cm.StaticSecret, cm.StaticToken)
	} else if credsType == CredentialsTypeShared {
		cfg.Credentials = credentials.NewSharedCredentials("", cm.SharedProfile)
	}
	if len(cm.Endpoint) > 0 {
		cfg.Endpoint = aws.String(cm.Endpoint)
	}
	if cm.Region == "" {
		cfg.Region = aws.String(cm.Region)
	} else {
		cfg.Region = aws.String(RegionUSWest1)
	}
	cfg.S3ForcePathStyle = aws.Bool(cm.PathStyleForce)
	return cfg
}

func (cm AWSConfigMore) NewSession() (*session.Session, error) {
	return session.NewSession(cm.Config())
	/*
		if cm.CredentialsType == CredentialsTypeEnvironment {
			return session.NewSession()
		}
		cfg := cm.Config(region)
		return session.NewSession(cfg)
	*/
}
