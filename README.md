# AWS Go (Helpers)

[![Build Status][build-status-svg]][build-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

AWSgo provices helper libraries for the AWS Go SDK:

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

 [build-status-svg]: https://github.com/grokify/awsgo/workflows/test/badge.svg
 [build-status-url]: https://github.com/grokify/awsgo/actions
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/awsgo
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/awsgo
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/awsgo
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/awsgo
 [loc-svg]: https://tokei.rs/b1/github/grokify/awsgo
 [repo-url]: https://github.com/grokify/awsgo
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/awsgo/blob/master/LICENSE
