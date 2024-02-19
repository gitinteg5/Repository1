package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/aws/aws-sdk-go-v2/service/configservice/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	configClient := configservice.NewFromConfig(cfg)

	
	channel, err := configClient.PutDeliveryChannel(context.TODO(), &configservice.PutDeliveryChannelInput{
		DeliveryChannel: &types.DeliveryChannel{
			ConfigSnapshotDeliveryProperties: nil,
			Name:                             nil,
			S3BucketName:                     nil,
			S3KeyPrefix:                      nil,
			S3KmsKeyArn:                      nil,
			SnsTopicARN:                      nil,
		},
	})
	if err != nil {
		log.Fatalf("unable to create delivery channel, %v", err)
	}
}