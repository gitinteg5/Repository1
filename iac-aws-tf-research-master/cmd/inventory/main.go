package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type S3Bucket struct {
	Name              string
	Region            string
	Policy            string
	IsPublicPolicy    bool
	ACL               S3BucketACL
	PublicAccessBlock S3BucketPublicAccessBlock
}

type S3BucketACL struct {
	OwnerDisplayName string
	OwnerID          string
	Grants           []s3types.Grant
}

type S3BucketPublicAccessBlock struct {
	BlockPublicAcls       bool
	BlockPublicPolicy     bool
	IgnorePublicAcls      bool
	RestrictPublicBuckets bool
}

func main() {
	customResolver := configureCustomResolver()

	var options []func(*config.LoadOptions) error

	if customResolver != nil {
		options = append(options, customResolver)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), options...)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	buckets, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("unable to list buckets, %v", err)
	}

	var s3Buckets []S3Bucket
	var limit = 5
	for idx, b := range buckets.Buckets {
		if idx >= limit {
			break
		}

		fmt.Printf("%s %v\n", *b.Name, b.CreationDate)

		region, _ := manager.GetBucketRegion(context.TODO(), s3Client, *b.Name)
		fmt.Printf("  - Region: %+v\n", region)

		policy, err := s3Client.GetBucketPolicy(context.TODO(), &s3.GetBucketPolicyInput{
			Bucket:              b.Name,
			ExpectedBucketOwner: nil,
		}, SetBucketRegion(region))
		if err != nil {
			var ae smithy.APIError
			if errors.As(err, &ae) && ae.ErrorCode() == "NoSuchBucketPolicy" {
				fmt.Printf("  - No policy \n")
			} else {
				log.Fatalf("unable to retrieve bucket policy, %v", err)
			}
		} else {
			fmt.Printf("  - Policy: %s\n", *policy.Policy)
		}

		acl, err := s3Client.GetBucketAcl(context.TODO(), &s3.GetBucketAclInput{
			Bucket:              b.Name,
			ExpectedBucketOwner: nil,
		}, SetBucketRegion(region))
		if err != nil {
			log.Fatalf("unable to retrieve bucket ACL, %v", err)
		} else {
			fmt.Printf("  - Owner: %+v - Grants: %+v\n", acl.Owner, acl.Grants)
		}

		policyStatus, err := s3Client.GetBucketPolicyStatus(context.TODO(), &s3.GetBucketPolicyStatusInput{
			Bucket:              b.Name,
			ExpectedBucketOwner: nil,
		}, SetBucketRegion(region))
		if err != nil {
			var ae smithy.APIError
			if errors.As(err, &ae) && ae.ErrorCode() == "NoSuchBucketPolicy" {
				fmt.Printf("  - No policy status\n")
			} else {
				log.Fatalf("unable to retrieve bucket policy status, %v", err)
			}
		} else {
			fmt.Printf("  - Policy status IsPublic: %+v\n", policyStatus.PolicyStatus.IsPublic)
		}

		publicAccessBlock, err := s3Client.GetPublicAccessBlock(context.TODO(), &s3.GetPublicAccessBlockInput{
			Bucket:              b.Name,
			ExpectedBucketOwner: nil,
		}, SetBucketRegion(region))
		if err != nil {
			var ae smithy.APIError
			if errors.As(err, &ae) && ae.ErrorCode() == "NoSuchPublicAccessBlockConfiguration" {
				fmt.Printf("  - No public access block\n")
			} else {
				log.Fatalf("unable to retrieve bucket public access block, %v", err)
			}
		} else {
			fmt.Printf("  - %+v\n", publicAccessBlock.PublicAccessBlockConfiguration)
		}

		fmt.Println()

		s3Bucket := S3Bucket{
			Name:              *b.Name,
			Region:            region,
		}

		if policy != nil {
			s3Bucket.Policy = *policy.Policy
		}

		if policyStatus != nil {
			s3Bucket.IsPublicPolicy = policyStatus.PolicyStatus.IsPublic
		}

		if publicAccessBlock != nil {
			s3Bucket.PublicAccessBlock = S3BucketPublicAccessBlock{
				BlockPublicAcls:       publicAccessBlock.PublicAccessBlockConfiguration.BlockPublicAcls,
				BlockPublicPolicy:     publicAccessBlock.PublicAccessBlockConfiguration.BlockPublicPolicy,
				IgnorePublicAcls:      publicAccessBlock.PublicAccessBlockConfiguration.IgnorePublicAcls,
				RestrictPublicBuckets: publicAccessBlock.PublicAccessBlockConfiguration.RestrictPublicBuckets,
			}
		}

		if acl != nil {
			s3Bucket.ACL = S3BucketACL{
				OwnerDisplayName: *acl.Owner.DisplayName,
				OwnerID:          *acl.Owner.ID,
				Grants:           acl.Grants,
			}
		}

		s3Buckets = append(s3Buckets, s3Bucket)
	}

	fmt.Println("JSON Output:")
	jsonData, _ := json.MarshalIndent(&s3Buckets, "", "  ")
	fmt.Printf("%s\n", jsonData)
}

func configureCustomResolver() config.LoadOptionsFunc {
	if os.Getenv("AWS_CUSTOM_ENDPOINT") == "" {
		return nil
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           os.Getenv("AWS_CUSTOM_ENDPOINT"),
			SigningRegion: "us-east-1",
		}, nil
	})

	return config.WithEndpointResolverWithOptions(customResolver)
}

func SetBucketRegion(region string) func(*s3.Options) {
	return func(o *s3.Options) {
		o.Region = region
	}
}
