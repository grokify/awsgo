package s3util

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/io/ioutil"
	"github.com/grokify/mogo/mime/mimeutil"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/http/httputilmore"
)

/*

https://pkg.go.dev/github.com/aws/aws-sdk-go

https://pkg.go.dev/github.com/aws/aws-sdk-go@v1.44.188/service/s3

https://questhenkart.medium.com/s3-image-uploads-via-aws-sdk-with-golang-63422857c548

https://www.backblaze.com/b2/docs/s3_compatible_api.html

https://help.backblaze.com/hc/en-us/articles/360047629713-Using-the-AWS-Go-SDK-with-B2

https://help.backblaze.com/hc/en-us/articles/360047425453-Getting-Started-with-the-S3-Compatible-API

Note: Buckets created prior to May 4th, 2020 are not S3 Compatible. If you do not have any S3 Compatible buckets, simply create a new bucket!

*/

// iterate all objects in a given S3 bucket and prefix, sum up objects' total size in bytes
// use: size, err := S3ObjectsSize("example-bucket-name", "/a/b/c")
// this works if (a) keypath is empty, (b) keypath is a prefix, (c) keypath is a full key.
func (cm *S3ClientMore) ObjectsKeys(bucket, keypath string) ([]string, error) {
	keys := []string{}
	// modified from https://stackoverflow.com/questions/67817746/get-s3-bucket-size-from-aws-go-sdk
	if cm.S3Client == nil {
		return keys, ErrS3ClientNotSet
	}
	output, err := cm.S3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(keypath),
	})

	if err != nil {
		return keys, fmt.Errorf("cannot ListObjectsV2 in (%s/%s): err (%s)", bucket, keypath, err.Error())
	}

	for _, object := range output.Contents {
		if object != nil && object.Key != nil {
			keys = append(keys, *object.Key)
		}
	}

	return keys, nil
}

func (cm *S3ClientMore) ObjectsKeysPrint(bucket, keypath string) ([]string, error) {
	keys, err := cm.ObjectsKeys(bucket, keypath)
	if err != nil {
		return keys, err
	}
	err = fmtutil.PrintJSON(keys)
	return keys, err
}

// ObjectsSize iterates all objects in a given S3 bucket and prefix to sum up objects' total size in bytes
// use: size, err := S3ObjectsSize("example-bucket-name", "/a/b/c")
func (cm *S3ClientMore) ObjectsSize(bucket, keypath string) (int64, error) {
	// modified from https://stackoverflow.com/questions/67817746/get-s3-bucket-size-from-aws-go-sdk
	if cm.S3Client == nil {
		return -1, ErrS3ClientNotSet
	}
	output, err := cm.S3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(keypath)})

	if err != nil {
		return -1, fmt.Errorf("cannot ListObjectsV2 in (%s/%s) err (%s)", bucket, keypath, err.Error())
	}

	var size int64
	for _, object := range output.Contents {
		if object != nil && object.Size != nil {
			size += *object.Size
		}
	}

	return size, nil
}

func (cm *S3ClientMore) ObjectStringPut(bucket, key, strBody string) (*s3.PutObjectOutput, error) {
	if cm.S3Client == nil {
		return nil, ErrS3ClientNotSet
	}
	return cm.S3Client.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader(strBody),
		Bucket: aws.String(bucket),
		Key:    aws.String(key)})
}

func (cm *S3ClientMore) ObjectFilePut(bucket, key, filename string) (*s3.PutObjectOutput, error) {
	if cm.S3Client == nil {
		return nil, ErrS3ClientNotSet
	}
	input, err := ObjectInputFileMore(bucket, key, filename)
	if err != nil {
		return nil, errorsutil.Wrap(err, "err in S3ClientMore.ObjectPutFile..UploadInputFile")
	}
	out, err := cm.S3Client.PutObject(input)
	if err != nil {
		return out, errorsutil.Wrapf(err, "failed to upload object %s/%s", bucket, key)
	}
	return out, nil
}

func (cm *S3ClientMore) ObjectFileUpload(bucket, key, filename string) (*s3manager.UploadOutput, error) {
	if cm.S3Client == nil {
		return nil, ErrS3ClientNotSet
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	ct, err := mimeutil.TypeByReadSeeker(f, true)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return cm.ObjectReaderUpload(bucket, key, ct, f)
}

func (cm *S3ClientMore) ObjectFileUploadOld(bucket, key, filename string) (*s3manager.UploadOutput, error) {
	if cm.S3Client == nil {
		return nil, ErrS3ClientNotSet
	} else if input, err := UploadInputFile(bucket, key, filename); err != nil {
		return nil, errorsutil.Wrap(err, "err in S3ClientMore.ObjectPutFile..UploadInputFile")
	} else {
		return cm.Uploader.Upload(input)
	}
}

func (cm *S3ClientMore) ObjectHTTPRequestPut(bucket, key string, sreq httpsimple.Request, sclient *httpsimple.Client) (*s3.PutObjectOutput, error) {
	err := checkBucketAndKey(bucket, key)
	if err != nil {
		return nil, err
	}
	if sclient == nil {
		sclient = &httpsimple.Client{}
	}
	if resp, err := sclient.Do(sreq); err != nil {
		return nil, err
	} else if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("http status code (%d)", resp.StatusCode)
	} else {
		return cm.ObjectHTTPResponsePut(bucket, key, resp)
	}
}

func (cm *S3ClientMore) ObjectHTTPRequestUpload(bucket, key string, sreq httpsimple.Request, sclient *httpsimple.Client) (*s3manager.UploadOutput, error) {
	err := checkBucketAndKey(bucket, key)
	if err != nil {
		return nil, err
	}
	if sclient == nil {
		sclient = &httpsimple.Client{}
	}
	resp, err := sclient.Do(sreq)
	if err != nil {
		return nil, err
	} else if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("http status code (%d)", resp.StatusCode)
	} else {
		return cm.ObjectReaderUpload(bucket, key, resp.Header.Get(httputilmore.HeaderContentType), resp.Body)
	}
}

func (cm *S3ClientMore) ObjectHTTPResponsePut(bucket, key string, resp *http.Response) (*s3.PutObjectOutput, error) {
	if cm.S3Client == nil {
		return nil, ErrS3ClientNotSet
	}
	err := checkBucketAndKey(bucket, key)
	if err != nil {
		return nil, err
	}
	input, err := ObjectInputHTTPResponseMore(bucket, key, resp)
	if err != nil {
		return nil, err
	}
	out, err := cm.S3Client.PutObject(input)
	if err != nil {
		return out, errorsutil.Wrapf(err, "failed to upload object %s/%s", bucket, key)
	}
	return out, nil
}

func ObjectInputHTTPResponseMore(bucket, key string, resp *http.Response) (*s3.PutObjectInput, error) {
	if resp == nil {
		return nil, errors.New("http.Response cannot be empty")
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("bad status code (%d)", resp.StatusCode)
	}
	// fmt.Printf("STATUS CODE [%d]\n", resp.StatusCode)
	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(resp.Header.Get(httputilmore.HeaderContentType)),
	}
	contentLength := resp.ContentLength
	if contentLength > 0 {
		rs, err := ioutil.ReaderToReadSeeker(resp.Body)
		if err != nil {
			return nil, err
		}
		input.Body = rs
		input.ContentLength = aws.Int64(contentLength)
	} else {
		ba, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if len(ba) == 0 {
			return nil, errors.New("body is empty")
		}
		input.Body = bytes.NewReader(ba)
		input.ContentLength = aws.Int64(int64(len(ba)))
	}
	return input, nil
}

func (cm *S3ClientMore) ObjectReaderUpload(bucket, key, contentType string, r io.Reader) (*s3manager.UploadOutput, error) {
	if cm.S3Client == nil {
		return nil, ErrS3ClientNotSet
	}
	input := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   r,
	}
	if strings.TrimSpace(contentType) != "" {
		input.ContentType = aws.String(contentType)
	}
	return cm.Uploader.Upload(input)
}

func ObjectInputFileMore(bucket, key, filename string) (*s3.PutObjectInput, error) {
	// Modified from https://help.backblaze.com/hc/en-us/articles/360047629713-Using-the-AWS-Go-SDK-with-B2
	input, err := ObjectInputFile(filename)
	if err != nil {
		return nil, err
	}
	input.Bucket = aws.String(bucket)
	input.Key = aws.String(key)
	return input, nil
}

func ObjectInputFile(filename string) (*s3.PutObjectInput, error) {
	// Modified from https://help.backblaze.com/hc/en-us/articles/360047629713-Using-the-AWS-Go-SDK-with-B2
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	var size = fileInfo.Size() // int64

	buffer := make([]byte, size)
	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	return &s3.PutObjectInput{
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}, nil
}

func UploadInputFile(bucket, key, filename string) (*s3manager.UploadInput, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	/*
		contentType, err := DetectContentTypeReadSeeker(f)
		if err != nil {
			return nil, err
		}*/

	fileInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}
	var size = fileInfo.Size() // int64

	buffer := make([]byte, size)
	_, err = f.Read(buffer)
	if err != nil {
		return nil, err
	}
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	return &s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        fileBytes,
		ContentType: aws.String(fileType),
	}, nil
}
