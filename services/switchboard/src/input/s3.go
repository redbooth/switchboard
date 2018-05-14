package input

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
)

type S3Conf struct {
	Endpoint string
	Region   string
	Profile  string
	Bucket   string
	Key      string
}

type S3 struct {
	conf    S3Conf
	errors  chan<- error
	readers chan<- io.ReadCloser
}

func NewS3(conf S3Conf, errors chan<- error, readers chan<- io.ReadCloser) *S3 {
	return &S3{conf, errors, readers}
}

func (input *S3) Read() {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String(input.conf.Endpoint),
			Region:   aws.String(input.conf.Region),
		},
		Profile: input.conf.Profile,
	})
	if err != nil {
		input.errors <- err
		log.Panicf("Unable to open AWS S3 download session %s/%s: %v", input.conf.Bucket, input.conf.Key, err)
	}

	svc := s3.New(sess)
	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(input.conf.Bucket),
		Key:    aws.String(input.conf.Key),
	})
	if err != nil {
		input.errors <- err
	} else {
		input.readers <- output.Body
	}
}
