package queue

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"strconv"
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
	urls, err := formatQueue(queuesOutput)
	if err != nil {
		return nil, err
	}

	return urls, err
}

func (h *Handler) GetQueue(ctx context.Context, queueURL string) (*Queue, error) {
	resQueue, err := h.Client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl: &queueURL,
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameApproximateNumberOfMessages,
			types.QueueAttributeNameVisibilityTimeout,
		},
	})
	if err != nil {
		return nil, err
	}

	messageCount, err := strconv.ParseUint(resQueue.Attributes[string(types.QueueAttributeNameApproximateNumberOfMessages)], 10, 32)
	if err != nil {
		return nil, err
	}

	messageTimeout, err := strconv.ParseUint(resQueue.Attributes[string(types.QueueAttributeNameVisibilityTimeout)], 10, 32)
	if err != nil {
		return nil, err
	}

	return &Queue{
		Name:           resolveNameFromUrl(queueURL),
		QueueARN:       resQueue.Attributes[string(types.QueueAttributeNameQueueArn)],
		MessageTimeout: uint32(messageTimeout),
		MessageCount:   uint32(messageCount),
	}, nil
}
