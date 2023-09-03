package rdsutil

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/grokify/mogo/pointer"
)

const (
	DBPGFAuroraMySQL57 = "aurora-mysql5.7"
	DBPGFMySQL80       = "mysql8.0"

// * aurora-mysql8.0
//
// * aurora-postgresql10
//
// * aurora-postgresql11
//
// * aurora-postgresql12
//
// * aurora-postgresql13
//
// * aurora-postgresql14
//
// * custom-oracle-ee-19
//
// * mariadb10.2
//
// * mariadb10.3
//
// * mariadb10.4
//
// * mariadb10.5
//
// * mariadb10.6
//
// * mysql5.7
//
// * mysql8.0
//
// * oracle-ee-19
//
// * oracle-ee-cdb-19
//
// * oracle-ee-cdb-21
//
// * oracle-se2-19
//
// * oracle-se2-cdb-19
//
// * oracle-se2-cdb-21
//
// * postgres10
//
// * postgres11
//
// * postgres12
//
// * postgres13
//
// * postgres14
//
// * sqlserver-ee-11.0
//
// * sqlserver-ee-12.0
//
// * sqlserver-ee-13.0
//
// * sqlserver-ee-14.0
//
// * sqlserver-ee-15.0
//
// * sqlserver-ex-11.0
//
// * sqlserver-ex-12.0
//
// * sqlserver-ex-13.0
//
// * sqlserver-ex-14.0
//
// * sqlserver-ex-15.0
//
// * sqlserver-se-11.0
//
// * sqlserver-se-12.0
//
// * sqlserver-se-13.0
//
// * sqlserver-se-14.0
//
// * sqlserver-se-15.0
//
// * sqlserver-web-11.0
//
// * sqlserver-web-12.0
//
// * sqlserver-web-13.0
//
// * sqlserver-web-14.0
//
// * sqlserver-web-15.0
)

var ErrPagesGT1000 = errors.New("more than 1000 pages returned")

type ParameterGroupService struct {
	rdsClient *rds.RDS
}

func (svc *ParameterGroupService) CreateDBClusterParameterGroup(
	dbClusterParameterGroupName string,
	dbParameterGroupFamily string,
	description string,
	opts *rds.CreateDBClusterParameterGroupInput) (*rds.CreateDBClusterParameterGroupOutput, error) {
	if opts == nil {
		opts = &rds.CreateDBClusterParameterGroupInput{}
	}
	if opts.DBClusterParameterGroupName == nil || strings.TrimSpace(*opts.DBClusterParameterGroupName) == "" {
		opts.DBClusterParameterGroupName = pointer.Pointer(strings.TrimSpace(dbClusterParameterGroupName))
	}
	if opts.DBParameterGroupFamily == nil || strings.TrimSpace(*opts.DBParameterGroupFamily) == "" {
		opts.DBParameterGroupFamily = pointer.Pointer(strings.TrimSpace(dbParameterGroupFamily))
	}
	if opts.Description == nil || strings.TrimSpace(*opts.Description) == "" {
		opts.Description = pointer.Pointer(strings.TrimSpace(description))
	}
	return svc.rdsClient.CreateDBClusterParameterGroup(opts)
}

func (svc *ParameterGroupService) DescribeDBParameters(dbParameterGroupName string, opts *rds.DescribeDBParametersInput) (Parameters, error) {
	params := []rds.Parameter{}
	if opts == nil {
		opts = &rds.DescribeDBParametersInput{}
	}
	if opts.DBParameterGroupName == nil || strings.TrimSpace(*opts.DBParameterGroupName) == "" {
		opts.DBParameterGroupName = pointer.Pointer(strings.TrimSpace(dbParameterGroupName))
	}
	pages := 0
	maxPages := 1000
	var marker string
	for {
		if marker != "" {
			opts.Marker = pointer.Pointer(marker)
		}
		res, err := svc.rdsClient.DescribeDBParameters(opts)
		if err != nil {
			return params, err
		}
		for _, p := range res.Parameters {
			params = append(params, *p)
		}
		if res.Marker == nil || strings.TrimSpace(*res.Marker) == "" {
			break
		} else {
			marker = *res.Marker
		}
		pages++
		if pages > maxPages {
			return params, ErrPagesGT1000
		}
	}
	return params, nil
}

// DescribeEngineDefaultParameters is a wrapper for https://docs.aws.amazon.com/sdk-for-go/api/service/rds/#RDS.DescribeEngineDefaultParameters
func (svc *ParameterGroupService) DescribeEngineDefaultParameters(dbParameterGroupFamily string, opts *rds.DescribeEngineDefaultParametersInput) (Parameters, error) {
	params := []rds.Parameter{}
	if opts == nil {
		opts = &rds.DescribeEngineDefaultParametersInput{}
	}
	if opts.DBParameterGroupFamily == nil || strings.TrimSpace(*opts.DBParameterGroupFamily) == "" {
		opts.DBParameterGroupFamily = pointer.Pointer(strings.TrimSpace(dbParameterGroupFamily))
	}
	if *opts.DBParameterGroupFamily == "" {
		return params, errors.New("dbParameterGroupFamily must be specified")
	}
	pages := 0
	maxPages := 1000
	var marker string
	for {
		if marker != "" {
			opts.Marker = pointer.Pointer(marker)
		}
		res, err := svc.rdsClient.DescribeEngineDefaultParameters(opts)
		if err != nil {
			return params, err
		}
		if res.EngineDefaults == nil {
			return params, errors.New("engineDefaults is nil")
		}
		for _, p := range res.EngineDefaults.Parameters {
			params = append(params, *p)
		}
		if res.EngineDefaults.Marker == nil || strings.TrimSpace(*res.EngineDefaults.Marker) == "" {
			break
		} else {
			marker = *res.EngineDefaults.Marker
		}
		pages++
		if pages > maxPages {
			return params, ErrPagesGT1000
		}
	}
	return params, nil
}

// ModifyDBParameterGroup is a wrapper for https://docs.aws.amazon.com/sdk-for-go/api/service/rds/#RDS.ModifyDBParameterGroup .
// If `dbParameterGroupName` or `params` are provided, they are used to override the value in `opts`.
func (svc *ParameterGroupService) ModifyDBParameterGroup(dbParameterGroupName string, params Parameters, opts *rds.ModifyDBParameterGroupInput) (*rds.DBParameterGroupNameMessage, error) {
	if opts != nil {
		opts = &rds.ModifyDBParameterGroupInput{}
	}
	dbParameterGroupName = strings.TrimSpace(dbParameterGroupName)
	if dbParameterGroupName != "" {
		opts.DBParameterGroupName = pointer.Pointer(dbParameterGroupName)
	}
	if len(params) > 0 {
		ptrs := params.ToPointers()
		opts.Parameters = ptrs
	}
	return svc.rdsClient.ModifyDBParameterGroup(opts)
}
