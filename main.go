package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	bucket, fileKey string
	client          *s3.Client
)

// PutFile - uploads the file to AWS.
func PutFile(c context.Context, client *s3.Client, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	fmt.Printf("Uploading: %v\n", aws.ToString(input.Key))
	return client.PutObject(c, input)
}

// GetFile - downloads the file from AWS.
func GetFile(ctx context.Context, client *s3.Client, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	fmt.Printf("Downloading: %v\n", aws.ToString(input.Key))

	return client.GetObject(ctx, input)
}

// SaveFile - saves downloaded file into current directory.
func SaveFile(file *s3.GetObjectOutput, key string) error {
	body, err := ioutil.ReadAll(file.Body)
	fileName := SplitKeyName(key)

	err = os.WriteFile(fileName, body, 0644)
	if err != nil {
		return err
	}
	return nil
}

// SplitKeyName - splits the given AWS file KEY and returns just only the file name.
func SplitKeyName(key string) string {
	splitName := strings.Split(key, "/")
	fileName := splitName[len(splitName)-1]
	return fileName
}

func main() {
	flag.StringVar(&bucket, "b", "", "The bucket to download/upload the file from/to")
	flag.StringVar(&fileKey, "f", "", "The file to download/upload")
	flag.Parse()

	if bucket == "" || fileKey == "" {
		fmt.Println("You must supply a bucket name (-b BUCKET) and file name (-f FILE)")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client = s3.NewFromConfig(cfg)

	//download the file
	dlInput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	}

	file, err := GetFile(context.TODO(), client, dlInput)
	if err != nil {
		log.Printf("GetFile error: %v", err)
		return
	}

	// do smth with file
	err = SaveFile(file, fileKey)
	if err != nil {
		log.Printf("WriteFile error: %v", err)
	}

	//upload the file
	ulFile, err := os.Open(SplitKeyName(fileKey))
	if err != nil {
		fmt.Println("Unable to open file " + SplitKeyName(fileKey))
		return
	}

	ulInput := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &fileKey,
		Body:   ulFile,
	}

	_, err = PutFile(context.TODO(), client, ulInput)
	if err != nil {
		log.Printf("Got error uploading file:%v\n", err)
		return
	}
}
