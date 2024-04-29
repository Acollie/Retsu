package queue

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Handler struct {
	Client *sqs.Client
	queues []*Queue
}
type HandlerI interface {
	RefreshQueue(ctx context.Context, queueARN string) error
	GetQueues(ctx context.Context) error
	GetQueue(ctx context.Context, queueARN string) (Queue, error)
}

type Queue struct {
	name           string
	queueARN       string
	messageCount   *uint32
	messageTimeout uint32
}

type QueueI interface {
	GetQueue(ctx context.Context) (Queue, error)
	RefreshQueue(ctx context.Context, queue Queue)
}
