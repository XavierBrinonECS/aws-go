package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

// Unmarshalling the json struct requires the member to start with capital letters.
type event struct {
	FileLocation string `json:"fileLocation"`
}

type response struct {
	FileLocation  string `json:"fileLocation"`
	FileAvailable bool   `json:"fileAvailable"`
}

func handler(evt event) (response, error) {
	// return fmt.Sprintf("File is %s", evt.fileLocation), nil

	if evt.FileLocation == "" {
		return response{}, fmt.Errorf("FileLocation is empty")
	}

	return response{
		FileLocation:  evt.FileLocation,
		FileAvailable: true,
	}, nil
}

func main() {
	lambda.Start(handler)
}
