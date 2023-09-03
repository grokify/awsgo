package rdsutil

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/grokify/awsgo/config"
)

type RDSClientMore struct {
	client *rds.RDS
}

func NewRDSClientMore(cm *config.AWSConfigMore) (*RDSClientMore, error) {
	if svc, err := NewRDSClient(cm); err != nil {
		return nil, err
	} else {
		return &RDSClientMore{client: svc}, nil
	}
}

func NewRDSClient(cm *config.AWSConfigMore) (*rds.RDS, error) {
	if cm == nil {
		return nil, errors.New("config.AWSConfigMore cannot be nil")
	}
	mySession, err := cm.NewSession()
	if err != nil {
		return nil, err
	}
	cfgs := []*aws.Config{}
	region := strings.TrimSpace(cm.Region)
	if region != "" {
		cfgs = append(cfgs, aws.NewConfig().WithRegion(region))
	}
	svc := rds.New(mySession, cfgs...)
	return svc, nil
}
