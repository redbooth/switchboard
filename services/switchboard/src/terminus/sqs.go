package terminus

import (
	"../header"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

type SQSConf struct {
	Endpoint string
	Region   string
	Profile  string
	QueueUrl string
}

type SQS struct {
	conf   SQSConf
	errors chan<- error
	svc    *sqs.SQS
}

func NewSQS(conf SQSConf, errors chan<- error) *SQS {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String(conf.Endpoint),
			Region:   aws.String(conf.Region),
		},
		Profile: conf.Profile,
	})
	if err != nil {
		log.Panicf("Unable to open AWS SQS session to queue %s: %v", conf.QueueUrl, err)
	}
	svc := sqs.New(sess)
	return &SQS{conf, errors, svc}
}

func (terminus *SQS) Terminate(h header.Header) {
	_, err := terminus.svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(h.String()),
		QueueUrl:    aws.String(terminus.conf.QueueUrl),
	})
	if err != nil {
		terminus.errors <- err
	}
}
