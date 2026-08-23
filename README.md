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
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://grokify.github.io/awsgo
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fawsgo
 [loc-svg]: https://tokei.rs/b1/github/grokify/awsgo
 [repo-url]: https://github.com/grokify/awsgo
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/awsgo/blob/main/LICENSE

AWSgo provides helper libraries for the AWS Go SDK:

* Dev Guide: https://aws.amazon.com/sdk-for-go/
* GitHub: https://github.com/aws/aws-sdk-go

## Installation

```
$ go get github.com/grokify/awsgo/...
```

## Configuration

AWS Basic Auth can be used using [`goauth`](https://github.com/grokify/goauth) with the following configuration:

```
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
``````
