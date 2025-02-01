package util

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadS3(uploadfile []byte, filename string) {

	var S3_BUCKET_NAME = os.Getenv("S3_BUCKET_NAME") //"bucket4image"
	var S3_BUCKET_REGION = os.Getenv("S3_BUCKET_REGION")
	var PROFILE = os.Getenv("PROFILE")
	/*
		var S3_BUCKET_NAME = "bucket4image"
		var S3_BUCKET_REGION = "ap-northeast-1"
		var PROFILE = "myregion"
	*/
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(S3_BUCKET_REGION),
		},
		Profile: PROFILE,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	//
	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(S3_BUCKET_NAME),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(uploadfile),
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}
