package output

import (
	"../header"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"log"
	"path"
)

type S3Conf struct {
	Endpoint  string
	Region    string
	Profile   string
	Bucket    string
	Directory string
	Extension string
}

type S3 struct {
	conf   S3Conf
	errors chan<- error
	header header.Header
	pr     *io.PipeReader
	pw     *io.PipeWriter
}

func NewS3(conf S3Conf, errors chan<- error, h header.Header) *S3 {
	pr, pw := io.Pipe()
	go func() {
		defer pr.Close()

		key := path.Join(conf.Directory, h.String())
		if len(conf.Extension) > 0 {
			key += "." + conf.Extension
		}

		sess, err := session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Endpoint: aws.String(conf.Endpoint),
				Region:   aws.String(conf.Region),
			},
			Profile: conf.Profile,
		})
		if err != nil {
			log.Panicf("Unable to open AWS S3 upload session %s/%s: %v", conf.Bucket, key, err)
		}

		uploader := s3manager.NewUploader(sess)
		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket:      aws.String(conf.Bucket),
			Key:         aws.String(key),
			ContentType: aws.String("application/octet-stream"),
			Body:        pr,
		})
		if err != nil {
			log.Printf("Failed to upload file %s/%s: %v\n", conf.Bucket, key, err)
		} else {
			log.Printf("Successfully uploaded file to %s\n", result.Location)
		}
	}()
	return &S3{conf, errors, h, pr, pw}
}

func (writer *S3) Write(b []byte) (n int, err error) {
	return writer.pw.Write(b)
}

func (writer *S3) Close() error {
	return writer.pw.Close()
}
