package awsx

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"queue/queue"
)

func NewHandler(config aws.Config) *queue.Handler {
	handler := &queue.Handler{
		Client: sqs.NewFromConfig(config),
	}

	return handler
}
