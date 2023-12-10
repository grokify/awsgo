// pricingutil provices some helpers for the AWS Pricing API: https://docs.aws.amazon.com/sdk-for-go/api/service/pricing/
package pricingutil

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/grokify/awsgo/config"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/type/stringsutil"
)

const (
	ServicesInputFormatVersionAWSV1 = "aws_v1"
	DescribeServicesMaxResults      = 100 // values > 100 will return `InvalidParameterException: 1 validation error detected: Value '10000' at 'maxResults' failed to satisfy constraint: Member must have value less than or equal to 100`
)

func NewPricingClient(cfg *config.AWSConfig) (*Client, error) {
	if cfg == nil {
		return nil, config.ErrAWSConfigCannotBeNil
	} else if ses, cfgs, err := cfg.ClientParams(); err != nil {
		return nil, err
	} else {
		return &Client{Pricing: pricing.New(ses, cfgs...)}, nil
	}
}

type Client struct {
	*pricing.Pricing
}

func (c *Client) Services() (Services, error) {
	svcs := Services{}
	params := pricing.DescribeServicesInput{
		FormatVersion: pointer.Pointer(ServicesInputFormatVersionAWSV1),
		MaxResults:    pointer.Pointer(int64(DescribeServicesMaxResults))}

	i := 0
	for {
		i++
		dsout, err := c.Pricing.DescribeServices(&params)
		if err != nil {
			return svcs, err
		} else if dsout == nil {
			return svcs, errors.New("return struct is nil")
		}
		for _, svc := range dsout.Services {
			if svc != nil {
				svcs = append(svcs, *svc)
			}
		}
		nextToken := strings.TrimSpace(pointer.Dereference(dsout.NextToken))
		if nextToken == "" {
			break
		} else if i > 10000 {
			return svcs, errors.New("over 10000 API next token iterations")
		}
		params.NextToken = pointer.Pointer(nextToken)
	}
	return svcs, nil
}

type Services []pricing.Service

func (svcs Services) ServiceCodes() ([]string, error) {
	codes := []string{}
	for _, svc := range svcs {
		if code := strings.TrimSpace(pointer.Dereference(svc.ServiceCode)); code != "" {
			codes = append(codes, code)
		}
	}
	return stringsutil.SliceCondenseSpace(codes, true, true), nil
}

func (c *Client) Products(svcCode string) (Products, error) {
	prods := Products{}
	if svcCode = strings.TrimSpace(svcCode); svcCode == "" {
		return prods, errors.New("service code cannot be empty")
	}

	params := &pricing.GetProductsInput{
		FormatVersion: pointer.Pointer(ServicesInputFormatVersionAWSV1),
		MaxResults:    pointer.Pointer(int64(DescribeServicesMaxResults)),
		ServiceCode:   pointer.Pointer(svcCode)}

	pageNum := 0
	err := c.Pricing.GetProductsPages(params,
		func(page *pricing.GetProductsOutput, lastPage bool) bool {
			if lastPage {
				return false
			} else if page != nil {
				prods = append(prods, page.PriceList...)
			}
			pageNum++
			return pageNum <= 100
		})

	return prods, err
}

type Products []aws.JSONValue
