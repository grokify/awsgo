package rdsutil

import (
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/grokify/awsgo/config"
)

type RDSClientMore struct {
	client    *rds.RDS
	Parameter *ParameterGroupService
}

func NewRDSClientMore(cfg *config.AWSConfig) (*RDSClientMore, error) {
	if cfg == nil {
		return nil, config.ErrAWSConfigCannotBeNil
	} else if svc, err := NewRDSClient(cfg); err != nil {
		return nil, err
	} else {
		return &RDSClientMore{
			client:    svc,
			Parameter: &ParameterGroupService{rdsClient: svc},
		}, nil
	}
}

/*
func NewRDSClientOld(cfg *config.AWSConfig) (*rds.RDS, error) {
	if cfg == nil {
		return nil, config.ErrAWSConfigCannotBeNil
	}
	mySession, err := cfg.NewSession()
	if err != nil {
		return nil, err
	}
	cfgs := []*aws.Config{}
	region := strings.TrimSpace(cfg.Region)
	if region != "" {
		cfgs = append(cfgs, aws.NewConfig().WithRegion(region))
	}
	svc := rds.New(mySession, cfgs...)
	return svc, nil
}
*/

func NewRDSClient(cfg *config.AWSConfig) (*rds.RDS, error) {
	if cfg == nil {
		return nil, config.ErrAWSConfigCannotBeNil
	} else if ses, cfgs, err := cfg.ClientParams(); err != nil {
		return nil, err
	} else {
		return rds.New(ses, cfgs...), nil
	}
}
