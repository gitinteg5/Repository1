terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

resource "aws_iam_account_password_policy" "example" {
  #Should be >= 14
  minimum_password_length        = 8
  #Should be >= 24
  password_reuse_prevention = 5
}