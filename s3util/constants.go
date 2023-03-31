package s3util

import (
	"fmt"
	"strings"
)

// Modified from https://help.backblaze.com/hc/en-us/articles/360047629713-Using-the-AWS-Go-SDK-with-B2

const (
	CredentialsTypeStatic = "static"

	DatacenterUSWest000 = "us-west-000" // https://github.com/peak/s5cmd/issues/269
	DatacenterUSWest002 = "us-west-002"

	backblazeS3EndpointFormat = `https://s3.%s.backblazeb2.com`
	BasePublicURL             = "https://f000.backblazeb2.com/file"
)

func B2Endpoint(region string) string {
	if strings.TrimSpace(region) == "" {
		region = DatacenterUSWest000
	}
	// https://s3.<region>.backblazeb2.com
	return fmt.Sprintf(backblazeS3EndpointFormat, strings.ToLower(strings.TrimSpace(region)))
}
