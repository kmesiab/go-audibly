# Existing S3 bucket for pre-transcribe
resource "aws_s3_bucket" "myp_pre_transcribe" {
  bucket = "myp-pre-transcribe"

  lifecycle {
    prevent_destroy = false
  }

  tags = {
    Name        = "myp"
    Environment = "dev"
  }
}

# Policy to allow Transcribe to access pre-transcribe bucket
resource "aws_s3_bucket_policy" "myp_pre_transcribe_policy" {
  bucket = aws_s3_bucket.myp_pre_transcribe.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect    = "Allow",
        Principal = { Service = "transcribe.amazonaws.com" },
        Action    = ["s3:GetObject", "s3:PutObject"],
        Resource  = "${aws_s3_bucket.myp_pre_transcribe.arn}/*"
      }
    ]
  })
}

# Existing S3 bucket for post-transcribe
resource "aws_s3_bucket" "myp_post_transcribe" {
  bucket = "myp-post-transcribe"

  lifecycle {
    prevent_destroy = false
    ignore_changes = [policy]
  }

  tags = {
    Name        = "myp"
    Environment = "dev"
  }
}

# Policy to allow Transcribe to access post-transcribe bucket
resource "aws_s3_bucket_policy" "myp_post_transcribe_policy" {
  bucket = aws_s3_bucket.myp_post_transcribe.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect    = "Allow",
        Principal = { Service = "transcribe.amazonaws.com" },
        Action    = ["s3:GetObject","s3:PutObject"],
        Resource  = "${aws_s3_bucket.myp_post_transcribe.arn}/*"
      }
    ]
  })
}

# Encryption for existing pre-transcribe bucket
resource "aws_kms_key" "s3key_pre" {
  description             = "This key is used to encrypt bucket objects for pre-transcribe"
  deletion_window_in_days = 10
}

resource "aws_s3_bucket_server_side_encryption_configuration" "pre_bucket" {
  bucket = aws_s3_bucket.myp_pre_transcribe.id

  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = aws_kms_key.s3key_pre.arn
      sse_algorithm     = "aws:kms"
    }
  }
}

# Encryption for new post-transcribe bucket
resource "aws_kms_key" "s3key_post" {
  description             = "This key is used to encrypt bucket objects for post-transcribe"
  deletion_window_in_days = 10
}

resource "aws_s3_bucket_server_side_encryption_configuration" "post_bucket" {
  bucket = aws_s3_bucket.myp_post_transcribe.id

  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = aws_kms_key.s3key_post.arn
      sse_algorithm     = "aws:kms"
    }
  }
}

# Versioning for pre-transcribe bucket
resource "aws_s3_bucket_versioning" "versioning_pre_bucket" {
  bucket = aws_s3_bucket.myp_pre_transcribe.id
  versioning_configuration {
    status = "Enabled"
  }
}

# Versioning for post-transcribe bucket
resource "aws_s3_bucket_versioning" "versioning_post_bucket" {
  bucket = aws_s3_bucket.myp_post_transcribe.id
  versioning_configuration {
    status = "Enabled"
  }
}
