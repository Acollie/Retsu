package queue

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"strings"
)

func resolveNameFromUrl(queueURL string) string {
	split := strings.Split(queueURL, "/")
	return split[len(split)-1]
}

func formatQueue(queuesOutputs []*sqs.ListQueuesOutput) ([]string, error) {
	var queues []string
	for _, output := range queuesOutputs {
		for _, url := range output.QueueUrls {
			queues = append(queues, resolveNameFromUrl(url))
		}
	}
	return queues, nil
}
