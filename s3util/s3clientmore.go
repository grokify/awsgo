package s3util

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/grokify/awsgo/config"
	"github.com/grokify/mogo/errors/errorsutil"
)

var (
	ErrAWSSessionNotSet = errors.New("aws session not set")
	ErrS3ClientNotSet   = errors.New("s3 client not set")
	ErrS3UploaderNotSet = errors.New("s3 uploader not set")
	ErrBucketNotSet     = errors.New("bucket not set")
	ErrKeyNotSet        = errors.New("key not set")
)

type S3ClientMore struct {
	Session    *session.Session
	S3Client   *s3.S3
	Downloader *s3manager.Downloader
	Uploader   *s3manager.Uploader
}

func NewS3Client(cm config.AWSConfigMore) (*s3.S3, *session.Session, error) {
	ses, err := cm.NewSession()
	if err != nil {
		return nil, ses, err
	}
	return s3.New(ses), ses, nil
}

func NewS3ClientMoreSession(sess *session.Session) (*S3ClientMore, error) {
	cm := &S3ClientMore{}
	return cm, cm.SetSession(sess)
}

func NewS3ClientMore(cm config.AWSConfigMore) (*S3ClientMore, error) {
	sess, err := cm.NewSession()
	if err != nil {
		return nil, err
	}
	return NewS3ClientMoreSession(sess)
}

func (cm *S3ClientMore) SetSession(sess *session.Session) error {
	if sess == nil {
		return ErrAWSSessionNotSet
	}
	cm.Session = sess
	cm.S3Client = s3.New(cm.Session)
	cm.Downloader = s3manager.NewDownloader(cm.Session) // Create a downloader with the session and default options
	cm.Uploader = s3manager.NewUploader(cm.Session)
	return nil
}

func (cm *S3ClientMore) CreateBucketSimple(bucket string) (*s3.CreateBucketOutput, error) {
	if cm.S3Client == nil {
		return nil, ErrS3ClientNotSet
	}
	return cm.S3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket), // Required
	})
}

// ObjectGetSimple returns an object given a bucket and key. For more complex
// usage, usage a custom `s3.GetObjectInput`.
func (cm *S3ClientMore) ObjectGetSimple(bucket, key string) (*s3.GetObjectOutput, error) {
	if cm.S3Client == nil {
		return nil, ErrS3ClientNotSet
	}
	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.PutObject
	out, err := cm.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return out, errorsutil.Wrap(err, "err in S3ClientMore.ObjectGetSimple..cm.S3Client.GetObject")
}

func (cm *S3ClientMore) ObjectGetToFile(bucket, key, filename string) (int64, error) {
	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/
	// Create a file to write the S3 Object contents to.
	f, err := os.Create(filename)
	if err != nil {
		return -1, fmt.Errorf("failed to create file %q, %v", filename, err)
	}
	defer f.Close()
	// Write the contents of S3 Object to the file
	return cm.Downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
}

func checkBucketAndKey(bucket, key string) error {
	if strings.TrimSpace(bucket) == "" {
		return ErrBucketNotSet
	} else if strings.TrimSpace(key) == "" {
		return ErrKeyNotSet
	}
	return nil
}
