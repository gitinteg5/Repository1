terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

resource "aws_iam_user" "example" {
  name = var.name
  path = var.path

  tags = {
    tag-key = "tag-value"
  }
}

resource "aws_iam_access_key" "example" {
  count = var.create_access_key ? 1 : 0
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