package main

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	// Additional imports
)

type extraction struct {
	item *dynamodb.GetItemOutput
	obj  []byte
}

func extract(ctx context.Context, partitionKey, sortKey string) (extraction, error) {
	var (
		ext    extraction
		objKey *string
		err    error
	)
	// Retrieve the item from DynamoDB
	ext.item, objKey, err = getItemFromDynamoDB(ctx, partitionKey, sortKey)
	if err != nil {
		return ext, err
	}
	// Retrieve the object associated with objKey
	if objKey != nil {
		ext.obj, err = getObjectFromS3(ctx, *objKey)
		if err != nil {
			return ext, err
		}
	}

	return ext, nil
}

func getItemFromDynamoDB(ctx context.Context, partitionKey, sortKey string) (*dynamodb.GetItemOutput, *string, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("YourTableName"), // Replace with your table name
		Key: map[string]types.AttributeValue{
			"PartitionKey": &types.AttributeValueMemberS{Value: partitionKey},
			"SortKey":      &types.AttributeValueMemberS{Value: sortKey},
		},
	}

	result, err := dynamoDBClient.GetItem(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	// Assuming the DynamoDB item has an 'S3Key' attribute
	s3Key, exists := result.Item["S3Key"]
	if !exists || s3Key == nil {
		// The S3 key does not exist or is nil
		return result, nil, nil
	}

	s3KeyValue, ok := s3Key.(*types.AttributeValueMemberS)
	if !ok {
		// The S3 key is not a string, handle accordingly
		return result, nil, nil
	}

	return result, &s3KeyValue.Value, nil
}

func getObjectFromS3(ctx context.Context, key string) ([]byte, error) {
	output, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String("BucketName"),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer output.Body.Close()

	return io.ReadAll(output.Body)
}
