resource "aws_cloudwatch_log_destination" "test" {
  name       = var.rName
  target_arn = aws_kinesis_stream.test.arn
  role_arn   = aws_iam_role.test.arn

{{- template "tags" . }}

  depends_on = [aws_iam_role_policy.test]
}

resource "aws_kinesis_stream" "test" {
  name        = var.rName
  shard_count = 1
}

data "aws_region" "current" {}

data "aws_iam_policy_document" "role" {
  statement {
    effect = "Allow"

    principals {
      type = "Service"

      identifiers = [
        "logs.${data.aws_region.current.region}.amazonaws.com",
      ]
    }

    actions = [
      "sts:AssumeRole",
    ]
  }
}

resource "aws_iam_role" "test" {
  name               = var.rName
  assume_role_policy = data.aws_iam_policy_document.role.json
}

data "aws_iam_policy_document" "policy" {
  statement {
    effect = "Allow"

    actions = [
      "kinesis:PutRecord",
    ]

    resources = [
      aws_kinesis_stream.test.arn,
    ]
  }

  statement {
    effect = "Allow"

    actions = [
      "iam:PassRole",
    ]

    resources = [
      aws_iam_role.test.arn,
    ]
  }
}

resource "aws_iam_role_policy" "test" {
  name   = var.rName
  role   = aws_iam_role.test.name
  policy = data.aws_iam_policy_document.policy.json
}
