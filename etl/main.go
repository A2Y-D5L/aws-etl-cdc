package main

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type key struct {
	partitionKey string
	sortKey      string
}

func unmarshalKey(body string) (key, error) {
	var k key
	err := json.Unmarshal([]byte(body), &k)
	if err != nil {
		// TODO: handle error
		return key{}, err
	}
	return k, nil
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	var wg sync.WaitGroup
	for _, rec := range sqsEvent.Records {
		wg.Add(1)
		go func(msg events.SQSMessage) {
			defer wg.Done()
			k, err := unmarshalKey(msg.Body)
			if err != nil {
				// TODO: handle error
				return
			}

			_, err = extract(ctx, k.partitionKey, k.sortKey)
			if err != nil {
				// TODO: handle error
				return
			}
		}(rec)
	}

	wg.Wait()
	return nil
}

func main() {
	lambda.Start(handler)
}
