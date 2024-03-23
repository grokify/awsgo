package config

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	aws2 "github.com/aws/aws-sdk-go-v2/aws"
	config2 "github.com/aws/aws-sdk-go-v2/config"
	credentials2 "github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/grokify/goauth"
	"github.com/grokify/mogo/pointer"
)

/*

https://docs.aws.amazon.com/sdk-for-go/api/aws/session/
https://docs.aws.amazon.com/sdk-for-go/api/service/s3/

https://pkg.go.dev/github.com/aws/aws-sdk-go-v2

*/

const (
	CredentialsTypeEnvironment = "environment" // uses AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN (optional)
	CredentialsTypeShared      = "shared"
	CredentialsTypeStatic      = "static"

	RegionUSEast1 = "us-east-1"
	RegionUSWest1 = "us-west-1"
	RegionDefault = RegionUSEast1 // from AWS: "If you don’t select a region, then us-east-1 will be used by default." https://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/setup-credentials.html
)

var ErrAWSConfigCannotBeNil = errors.New("config.AWSConfig cannot be nil")

// AWSConfig handles credentials from https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
type AWSConfig struct {
	CredentialsType string
	SharedProfile   string // `.aws/credentials`
	StaticID        string
	StaticSecret    string
	StaticToken     string // optional per https://docs.aws.amazon.com/sdk-for-go/api/aws/session/
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

func (cfg AWSConfig) Config() *aws.Config {
	ac := &aws.Config{}
	credsType := strings.ToLower(strings.TrimSpace(cfg.CredentialsType))

	if credsType == CredentialsTypeStatic {
		ac.Credentials = credentials.NewStaticCredentials(cfg.StaticID, cfg.StaticSecret, cfg.StaticToken)
	} else if credsType == CredentialsTypeShared {
		ac.Credentials = credentials.NewSharedCredentials("", cfg.SharedProfile)
	}
	if len(cfg.Endpoint) > 0 {
		ac.Endpoint = aws.String(cfg.Endpoint)
	}
	/*
		if cm.Region == "" {
			ac.Region = aws.String(cm.Region)
		} else {
			ac.Region = aws.String(RegionUSWest1)
		}
	*/
	ac.Region = pointer.Pointer(cfg.RegionOrDefault(RegionDefault))
	ac.S3ForcePathStyle = aws.Bool(cfg.PathStyleForce)
	return ac
}

func (cfg AWSConfig) ConfigV2(ctx context.Context) (aws2.Config, error) {
	// https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/
	if ac, err := config2.LoadDefaultConfig(
		ctx,
		config2.WithCredentialsProvider(
			credentials2.NewStaticCredentialsProvider(
				cfg.StaticID, cfg.StaticSecret, cfg.StaticToken)),
	); err != nil {
		return ac, err
	} else {
		ac.Region = cfg.RegionOrDefault(RegionDefault)
		return ac, nil
	}
}

func (cfg AWSConfig) RegionOrDefault(def string) string {
	if region := strings.TrimSpace(cfg.Region); region == "" {
		return def
	} else {
		return region
	}
}

func (cfg AWSConfig) NewSession() (*session.Session, error) {
	return session.NewSession(cfg.Config())
	/*
		if cm.CredentialsType == CredentialsTypeEnvironment {
			return session.NewSession()
		}
		cfg := cm.Config(region)
		return session.NewSession(cfg)
	*/
}

// ClientParams returns the params used to set up a service.
func (cfg AWSConfig) ClientParams() (client.ConfigProvider, []*aws.Config, error) {
	ses, err := cfg.NewSession()
	if err != nil {
		return nil, []*aws.Config{}, err
	}
	cfgs := []*aws.Config{cfg.Config()}
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
