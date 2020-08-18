package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Unmarshalling the json struct requires the member to start with capital letters.
type event struct {
	FileLocation string `json:"fileLocation"`
	BucketName   string `json:"bucketName"`
}

type response struct {
	FileLocation  string `json:"fileLocation"`
	BucketName    string `json:"bucketName"`
	FileAvailable bool   `json:"fileAvailable"`
}

func handler(ctx context.Context, evt event) (response, error) {
	// return fmt.Sprintf("File is %s", evt.fileLocation), nil

	if evt.BucketName == "" {
		return response{}, fmt.Errorf("BucketName is empty")
	}

	if evt.FileLocation == "" {
		return response{}, fmt.Errorf("FileLocation is empty")
	}

	s3Service := s3.New(session.New())
	input := &s3.HeadObjectInput{
		Bucket: aws.String(evt.BucketName),
		Key:    aws.String(evt.FileLocation),
	}

	res, err := s3Service.HeadObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}
	fmt.Println("res")
	fmt.Println(res.GoString())
	fmt.Println(res.LastModified)
	fmt.Println(res.LastModified != nil)

	return response{
		FileLocation:  evt.FileLocation,
		BucketName:    evt.BucketName,
		FileAvailable: res.LastModified != nil,
	}, nil
}

func main() {
	lambda.Start(handler)
}
