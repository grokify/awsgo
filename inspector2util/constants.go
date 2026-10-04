package inspector2util

import (
	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
)

const (
	sepFilepathVersion = "@"
	sepJoinDispolay    = ", "

	yesnoYes = "Yes"
	yesnoNo  = "No"

	filenamePomProperties = "pom.properties"

	// Resource table field-name constants.
	ImageHash           = "image_hash"
	ImageRepositoryName = "image_repo_name"
	ImageTags           = "image_tags"

	ResourceTypeAwsEc2Instance       = types.ResourceTypeAwsEc2Instance
	ResourceTypeAwsEcrContainerImage = types.ResourceTypeAwsEcrContainerImage
	ResourceTypeAwsEcrRepository     = types.ResourceTypeAwsEcrRepository
	ResourceTypeAwsLambdaFunction    = types.ResourceTypeAwsLambdaFunction
)
