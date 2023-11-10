resource "aws_iam_user" "my_user" {
  name = "go-audibly"
  tags = {
    "Name" = "app user"
  }
}

resource "aws_iam_access_key" "my_access_key" {
  user = aws_iam_user.my_user.name
}

output "aws_access_key_id" {
  value = aws_iam_access_key.my_access_key.id
}

output "aws_secret_access_key" {
  value     = aws_iam_access_key.my_access_key.secret
  sensitive = true
}

locals {
  secret_path = "aws_user_credentials.txt"
}

resource "local_file" "secret_key" {
  content  = aws_iam_access_key.my_access_key.secret
  filename = local.secret_path
}

output "secret_file_path" {
  value     = local.secret_path
  sensitive = true
}

resource "aws_iam_policy" "s3_full_access_policy" {
  name        = "S3FullAccessPolicy"
  description = "S3 full permissions for go-audibly to upload audio files"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "s3:*",
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_user_policy_attachment" "user_policy_attachment" {
  user       = aws_iam_user.my_user.name
  policy_arn = aws_iam_policy.s3_full_access_policy.arn
}

resource "aws_iam_policy" "kms_generate_data_key" {
  name        = "kms_generate_data_key"
  description = "Allows KMS GenerateDataKey operations"

  policy = jsonencode({
    Version   = "2012-10-17",
    Statement = [
      {
        Action   = "kms:GenerateDataKey",
        Effect   = "Allow",
        Resource = "arn:aws:kms:us-west-2:462498369025:key/*"
      }
    ]
  })
}

resource "aws_iam_policy_attachment" "attach_kms_generate_data_key" {
  name       = "attach_kms_generate_data_key"
  users      = [aws_iam_user.my_user.name]
  policy_arn = aws_iam_policy.kms_generate_data_key.arn
}

resource "aws_iam_role" "transcribe_service_role" {
  name = "TranscribeServiceRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Service = "transcribe.amazonaws.com"
        },
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy" "transcribe_s3_access" {
  name = "TranscribeS3Access"
  role = aws_iam_role.transcribe_service_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "s3:*",
        ],
        Resource = [
          "${aws_s3_bucket.audio_bucket.arn}/*",
          "${aws_s3_bucket.transcript_bucket.arn}/*"
        ]
      },
      {
        Effect = "Allow",
        Action = [
          "s3:*"
        ],
        Resource = [
          "${aws_s3_bucket.transcript_bucket.arn}/*"
        ]
      }
    ]
  })
}

resource "aws_s3_bucket_policy" "transcript_bucket_policy" {
  bucket = aws_s3_bucket.transcript_bucket.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect    = "Allow",
        Principal = { Service = "transcribe.amazonaws.com" },
        Action    = "s3:*",
        Resource  = "${aws_s3_bucket.transcript_bucket.arn}/*"
      },
    ]
  })
}

resource "aws_iam_role_policy" "transcribe_kms_access" {
  name = "TranscribeKMSAccess"
  role = aws_iam_role.transcribe_service_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = "kms:Decrypt",
        Resource = [
          aws_kms_key.s3key_pre.arn,
          aws_kms_key.s3key_post.arn
        ]
      }
    ]
  })
}

