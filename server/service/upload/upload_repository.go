package upload

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// const (
// 	region     = "ap-south-1"
// 	accessKey  = "AKIA5PARVWEFTXOUKFSQ"
// 	secretKey  = "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"
// 	bucketName = "aws-bucket-media-storage"
// )

var uploder *s3manager.Uploader

type AwsService struct {
	region     string
	accessKey  string
	secretKey  string
	bucketName string
}

func NewAwsService(region, accessKey, secretKey, bucketName string) *AwsService {
	return &AwsService{
		region:     region,
		accessKey:  accessKey,
		secretKey:  secretKey,
		bucketName: bucketName,
	}
}

func (a *AwsService) InitializeAwsSerive(region, accessKey, secretKey string) {
	awsSession, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentials(
				accessKey, secretKey, "",
			),
		},
	})

	if err != nil {
		log.Println(err)
	}

	uploder = s3manager.NewUploader(awsSession)
}

func (a *AwsService) UploaFile(files []*multipart.FileHeader) ([]string, error) {

	var errorFiles []string
	var uploadedURLs []string

	for _, file := range files {
		fileHeader := file

		f, err := fileHeader.Open()
		if err != nil {
			errorFiles = append(errorFiles, fmt.Sprintf("error opening file %s: %s", fileHeader.Filename, err.Error()))
			continue
		}
		defer f.Close()

		uploadedURL, err := a.saveFile(f, fileHeader)
		if err != nil {
			errorFiles = append(errorFiles, fmt.Sprintf("error opening file %s: %s", fileHeader.Filename, err.Error()))
		} else {
			uploadedURLs = append(uploadedURLs, uploadedURL)
		}
	}

	if len(errorFiles) > 0 {
		return errorFiles, errors.New("error occured while image uploading")
	}

	return uploadedURLs, nil
}

func (a *AwsService) saveFile(fileReader io.Reader, fileHeader *multipart.FileHeader) (string, error) {

	// upload the file
	_, err := uploder.Upload(&s3manager.UploadInput{
		Bucket: aws.String(a.bucketName),
		Key:    aws.String(fileHeader.Filename),
		Body:   fileReader,
	})

	if err != nil {
		return "", err
	}

	// get the uploaded file URL
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", a.bucketName, fileHeader.Filename)
	return url, nil
}
