
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
  value = aws_iam_access_key.my_access_key.secret
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


