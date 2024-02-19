terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

resource "aws_iam_user" "example" {
  name = "example_user"
  path = "/foo/"

  tags = {
    tag-key = "tag-value"
  }
}

resource "aws_iam_access_key" "example" {
  user = aws_iam_user.example.name
}

resource "aws_iam_user_policy" "example_ro" {
  name = "test"
  user = aws_iam_user.example.name

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:Describe*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}