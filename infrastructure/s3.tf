resource "aws_s3_bucket" "myp_pre_transcode" {
  bucket = "myp-pre-transcode"
  tags = {
    Name        = "myp"
    Environment = "dev"
  }
}

resource "aws_s3_bucket_versioning" "versioning_example" {
  bucket = aws_s3_bucket.myp_pre_transcode.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_policy" "bucket_policy" {
  bucket = aws_s3_bucket.myp_pre_transcode.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          AWS = aws_iam_user.my_user.arn
        }
        Action = ["s3:PutObject", "s3:GetObject"]
        Resource = ["${aws_s3_bucket.myp_pre_transcode.arn}/*"]
      }
    ]
  })
}
