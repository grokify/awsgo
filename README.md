# AWS Go (Helpers)

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/grokify/awsgo/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/grokify/awsgo/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/grokify/awsgo/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/grokify/awsgo/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/grokify/awsgo/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/grokify/awsgo/actions/workflows/go-sast-codeql.yaml
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/awsgo
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/awsgo
 [viz-svg]: https://img.shields.io/badge/visualization-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fawsgo
 [loc-svg]: https://tokei.rs/b1/github/grokify/awsgo
 [repo-url]: https://github.com/grokify/awsgo
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/awsgo/blob/main/LICENSE

`awsgo` is a collection of Go helper libraries for the [AWS SDK for Go](https://aws.amazon.com/sdk-for-go/) ([v1](https://github.com/aws/aws-sdk-go) and [v2](https://github.com/aws/aws-sdk-go-v2)), covering configuration/auth, several AWS service clients, and AWS/Azure region reference data. It follows a library-first design: reusable packages with thin CLI adapters on top.

## Installation

```bash
go get github.com/grokify/awsgo/...
```

## Packages

| Package | Description |
|---|---|
| [`config`](config) | AWS SDK v1/v2 configuration, sessions, and credentials, with [`goauth`](https://github.com/grokify/goauth) integration |
| [`costexplorerutil`](costexplorerutil) | Helpers for the [AWS Cost Explorer](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/costexplorer) API |
| [`dynamodbutil`](dynamodbutil) | A simple key-value store client built on DynamoDB |
| [`iamutil`](iamutil) | Helpers for IAM policies, roles, users, and services |
| [`inspector2util`](inspector2util) | Reads AWS Inspector2 finding exports and provides aggregation, filtering, and streaming JSON parsing |
| [`pricingutil`](pricingutil) | Helpers for the [AWS Pricing](https://docs.aws.amazon.com/sdk-for-go/api/service/pricing/) API |
| [`rdsutil`](rdsutil) | RDS parameter group helpers: read, compare engine defaults, diff parameter sets |
| [`regions`](regions) | AWS and Azure region/location reference data (UN/LOCODE, ISO 3166 country and subdivision codes) |
| [`s3util`](s3util) | S3 object and key helpers |
| [`secretsmanagerutil`](secretsmanagerutil) | Helpers for AWS Secrets Manager |
| [`textractutil`](textractutil) | Helpers for AWS Textract, including an [`ocrutil`](textractutil/ocrutil) subpackage |
| [`trustedadvisorutil`](trustedadvisorutil) | Helpers for the [AWS Trusted Advisor](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/trustedadvisor) API |

Full API documentation is available on [pkg.go.dev][docs-godoc-url].

## CLI

[`cmd/awsgo`](cmd/awsgo) is a growing command-line toolbox built on these packages. It hasn't been tagged in a release yet, so build it from a clone for now:

```bash
git clone https://github.com/grokify/awsgo.git
cd awsgo
go build -o awsgo ./cmd/awsgo

./awsgo rds param-read --file rdsutil/testdata/db-perf-group_mysql8.json
```

Run `awsgo --help` for the current list of commands. Several packages also ship standalone example programs under their own `cmd/` subdirectories (e.g. `costexplorerutil/cmd`, `pricingutil/cmd`, `regions/cmd`, `trustedadvisorutil/cmd`).

## Configuration

AWS Basic Auth can be used via [`goauth`](https://github.com/grokify/goauth) with the following configuration:

```json
{
	"credentials": {
		"AWS": {
			"service": "aws",
			"type": "basic",
			"basic": {
				"username": "my_username",
				"password": "my_password"
			}
		}
	}
}
```

## License

[MIT][license-url]
