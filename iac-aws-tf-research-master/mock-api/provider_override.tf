provider "aws" {
  #Some endpoints can be replaced with env vars. See https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/custom-service-endpoints
  endpoints {
    ec2 = "http://localhost:3000" #Alternatively env TF_AWS_S3_ENDPOINT
    s3  = "http://localhost:3000" #Alternatively env TF_AWS_S3_ENDPOINT
    sts = "http://localhost:3000" #Alternatively env TF_AWS_STS_ENDPOINT
  }
  region = "us-east-1"
  skip_requesting_account_id = true
}
