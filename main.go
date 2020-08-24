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
type fileDescriptor struct {
	FileLocation string `json:"fileLocation"`
	BucketName   string `json:"bucketName"`
}

type event struct {
	Files []fileDescriptor `json:"files"`
}

type response struct {
	AllFilesAvailable bool `json:"allFilesAvailable"`
}

func isFileAvailable(file fileDescriptor) (bool, error) {
	if file.BucketName == "" {
		return false, fmt.Errorf("BucketName is empty")
	}

	if file.FileLocation == "" {
		return false, fmt.Errorf("FileLocation is empty")
	}

	s3Service := s3.New(session.New())
	input := &s3.HeadObjectInput{
		Bucket: aws.String(file.BucketName),
		Key:    aws.String(file.FileLocation),
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

	return res.LastModified != nil, nil
}

func handler(ctx context.Context, evt event) (response, error) {

	// Just for log purpose
	areAvailable := make([]bool, len(evt.Files))
	for i, fileDescr := range evt.Files {
		areAvailable[i], _ = isFileAvailable(fileDescr)
	}

	allAvailable := true
	for _, fileDescr := range evt.Files {
		allAvailable, _ = isFileAvailable(fileDescr)
		if !allAvailable {
			break
		}
	}

	fmt.Println("areAvailable", areAvailable)

	return response{
		AllFilesAvailable: allAvailable,
	}, nil
}

func main() {
	lambda.Start(handler)
}
