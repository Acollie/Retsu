package queue

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func (h *Handler) GetQueues(ctx context.Context) error {
	paginator := sqs.NewListQueuesPaginator(h.Client, &sqs.ListQueuesInput{})
	var queuesOutput []*sqs.ListQueuesOutput
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		queuesOutput = append(queuesOutput, output)
	}
	queues, err := formatQueue(ctx, h.Client, queuesOutput)
	if err != nil {
		return err
	}
	h.queues = queues
	return nil
}

func formatQueue(ctx context.Context, client *sqs.Client, queuesOutputs []*sqs.ListQueuesOutput) ([]*Queue, error) {
	var queues []*Queue
	for _, output := range queuesOutputs {
		for _, url := range output.QueueUrls {
			queue, err := getQueue(ctx, client, url)
			if err != nil {
				return nil, err
			}
			queues = append(queues, queue)
		}

	}
	return queues, nil

}

func getQueue(ctx context.Context, client *sqs.Client, queueURL string) (*Queue, error) {
	return nil, nil
}
