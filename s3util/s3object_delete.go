package s3util

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// MSSToObjectIdentifiers accepts a `map[string]string` of object info where
// keys are `ObjectIdentifier.Key` and the values are `ObjectIdentifier.VersionId`.`
func MSSToObjectIdentifiers(m map[string]string) []*s3.ObjectIdentifier {
	var ois []*s3.ObjectIdentifier
	for k, v := range m {
		oi := s3.ObjectIdentifier{
			Key: aws.String(k)}
		if v != "" {
			oi.VersionId = aws.String(v)
		}
		ois = append(ois, &oi)
	}
	return ois
}

func DeleteObjectsInputFromKeys(bucket string, keys []string, quiet bool) *s3.DeleteObjectsInput {
	input := &s3.DeleteObjectsInput{
		Bucket: aws.String(bucket)}
	del := &s3.Delete{}
	if quiet {
		del.Quiet = aws.Bool(quiet)
	}
	m := map[string]string{}
	for _, k := range keys {
		m[k] = ""
	}
	del.Objects = MSSToObjectIdentifiers(m)
	input.Delete = del
	return input
}

func (cm *S3ClientMore) DeleteObjectsByKeys(bucket string, keys []string, quiet bool) (*s3.DeleteObjectsOutput, error) {
	if cm.S3Client == nil {
		return nil, ErrS3ClientNotSet
	}
	return cm.S3Client.DeleteObjects(DeleteObjectsInputFromKeys(bucket, keys, quiet))
}
