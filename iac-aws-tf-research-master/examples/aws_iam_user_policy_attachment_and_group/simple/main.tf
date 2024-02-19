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

resource "aws_iam_policy" "policy" {
  name        = "too-much-permissions"
  description = "A policy allowing *:* violating CIS 1.16"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_user_policy_attachment" "test-attach" {
  user       = aws_iam_user.example.name
  policy_arn = aws_iam_policy.policy.arn
}

resource "aws_iam_group" "example" {
  name = "example"
  path = "/users/"
}

resource "aws_iam_group_membership" "team" {
  name = "test-group-membership"

  users = [
    aws_iam_user.example.name
  ]

  group = aws_iam_group.example.name
}

resource "aws_iam_group_policy" "group_policy" {
  name  = "some-test-policy"
  group = aws_iam_group.example.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "ec2:Describe*",
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}

resource "aws_iam_group_policy_attachment" "test-attach" {
  group      = aws_iam_group.example.name
  policy_arn = aws_iam_policy.policy.arn
}