package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/grokify/goauth"
	"github.com/grokify/mogo/pointer"
)

/*

https://docs.aws.amazon.com/sdk-for-go/api/aws/session/
https://docs.aws.amazon.com/sdk-for-go/api/service/s3/

*/

const (
	CredentialsTypeEnvironment = "environment" // uses AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN (optional)
	CredentialsTypeShared      = "shared"
	CredentialsTypeStatic      = "static"

	RegionUSEast1 = "us-east-1"
	RegionUSWest1 = "us-west-1"
)

var ErrAWSConfigCannotBeNil = errors.New("config.AWSConfig cannot be nil")

// AWSConfig handles credentials from https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
type AWSConfig struct {
	CredentialsType string
	SharedProfile   string // `.aws/credentials`
	StaticID        string
	StaticSecret    string
	StaticToken     string
	Endpoint        string
	Region          string
	PathStyleForce  bool
}

func NewAWSConfigStatic(region, accessKeyID, accessKeySecret string) *AWSConfig {
	return &AWSConfig{
		CredentialsType: CredentialsTypeStatic,
		StaticID:        accessKeyID,
		StaticSecret:    accessKeySecret,
		Region:          region,
	}
}

func AWSConfigCredentialsSetFile(filename, key string) (*AWSConfig, error) {
	if cs, err := goauth.ReadFileCredentialsSet(filename, false); err != nil {
		return nil, err
	} else {
		return AWSConfigCredentialsSet(cs, key)
	}
}

func AWSConfigCredentialsSet(set *goauth.CredentialsSet, key string) (*AWSConfig, error) {
	if set == nil {
		return nil, errors.New("credentials set cannot be nil")
	} else if creds, ok := set.Credentials[key]; !ok {
		return nil, fmt.Errorf("credentials key (%s) not found", key)
	} else {
		return AWSConfigMoreCredentialsBasic(&creds)
	}
}

func AWSConfigMoreCredentialsBasic(creds *goauth.Credentials) (*AWSConfig, error) {
	if creds == nil {
		return nil, errors.New("goauth credentials cannot be nil")
	} else if creds.Type == goauth.TypeBasic && creds.Basic != nil {
		return NewAWSConfigStatic(
			"",
			creds.Basic.Username,
			creds.Basic.Password), nil
	} else {
		return nil, errors.New("creds type not supported")
	}
}

func (cm AWSConfig) Config() *aws.Config {
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
	/*
		if cm.Region == "" {
			cfg.Region = aws.String(cm.Region)
		} else {
			cfg.Region = aws.String(RegionUSWest1)
		}
	*/
	cfg.Region = pointer.Pointer(cm.RegionOrDefault(RegionUSEast1))
	cfg.S3ForcePathStyle = aws.Bool(cm.PathStyleForce)
	return cfg
}

func (cm AWSConfig) RegionOrDefault(def string) string {
	if region := strings.TrimSpace(cm.Region); region == "" {
		return def
	} else {
		return region
	}
}

func (cm AWSConfig) NewSession() (*session.Session, error) {
	return session.NewSession(cm.Config())
	/*
		if cm.CredentialsType == CredentialsTypeEnvironment {
			return session.NewSession()
		}
		cfg := cm.Config(region)
		return session.NewSession(cfg)
	*/
}

// ClientParams returns the params used to set up a service.
func (acf AWSConfig) ClientParams() (client.ConfigProvider, []*aws.Config, error) {
	ses, err := acf.NewSession()
	if err != nil {
		return nil, []*aws.Config{}, err
	}
	cfgs := []*aws.Config{acf.Config()}
	/*
		cfgs := []*aws.Config{}
		region := strings.TrimSpace(cfg.Region)
		if region != "" {
			cfgs = append(cfgs, aws.NewConfig().WithRegion(region))
		}
	*/
	return ses, cfgs, nil
}

// p client.ConfigProvider, cfgs ...*aws.Config) *Pricing
