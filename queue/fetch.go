package queue

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"strings"
)

func (h *Handler) Scan(ctx context.Context) ([]string, error) {
	paginator := sqs.NewListQueuesPaginator(h.Client, &sqs.ListQueuesInput{})
	var queuesOutput []*sqs.ListQueuesOutput
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		queuesOutput = append(queuesOutput, output)
	}
	urls, err := formatQueue(ctx, h.Client, queuesOutput)
	if err != nil {
		return nil, err
	}

	return urls, err
}

func formatQueue(ctx context.Context, client *sqs.Client, queuesOutputs []*sqs.ListQueuesOutput) ([]string, error) {
	var queues []string
	for _, output := range queuesOutputs {
		for _, url := range output.QueueUrls {
			queues = append(queues, resolveNameFromUrl(url))
		}

	}
	return queues, nil

}

func (h *Handler) GetQueue(ctx context.Context, queueURL string) (*Queue, error) {
	resQueue, err := h.Client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl: &queueURL,
	})
	if err != nil {
		return nil, err
	}
	return &Queue{
		Name:     resolveNameFromUrl(queueURL),
		QueueARN: resQueue.Attributes["QueueArn"],
	}, nil

}

func resolveNameFromUrl(queueURL string) string {
	split := strings.Split(queueURL, "/")
	return split[len(split)-1]
}
