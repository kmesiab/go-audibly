# Existing S3 bucket for pre-transcode
resource "aws_s3_bucket" "myp_pre_transcode" {
  bucket = "myp-pre-transcode"
  tags = {
    Name        = "myp"
    Environment = "dev"
  }
}

# Policy to allow Transcribe to access pre-transcode bucket
resource "aws_s3_bucket_policy" "myp_pre_transcode_policy" {
  bucket = aws_s3_bucket.myp_pre_transcode.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect    = "Allow",
        Principal = { Service = "transcribe.amazonaws.com" },
        Action    = "s3:GetObject",
        Resource  = "${aws_s3_bucket.myp_pre_transcode.arn}/*"
      }
    ]
  })
}

# Existing S3 bucket for post-transcode
resource "aws_s3_bucket" "myp_post_transcode" {
  bucket = "myp-post-transcode"
  tags = {
    Name        = "myp"
    Environment = "dev"
  }
}

# Policy to allow Transcribe to access post-transcode bucket
resource "aws_s3_bucket_policy" "myp_post_transcode_policy" {
  bucket = aws_s3_bucket.myp_post_transcode.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect    = "Allow",
        Principal = { Service = "transcribe.amazonaws.com" },
        Action    = "s3:GetObject",
        Resource  = "${aws_s3_bucket.myp_post_transcode.arn}/*"
      }
    ]
  })
}

# Encryption for existing pre-transcode bucket
resource "aws_kms_key" "s3key_pre" {
  description             = "This key is used to encrypt bucket objects for pre-transcode"
  deletion_window_in_days = 10
}

resource "aws_s3_bucket_server_side_encryption_configuration" "pre_example" {
  bucket = aws_s3_bucket.myp_pre_transcode.id

  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = aws_kms_key.s3key_pre.arn
      sse_algorithm     = "aws:kms"
    }
  }
}

# Encryption for new post-transcode bucket
resource "aws_kms_key" "s3key_post" {
  description             = "This key is used to encrypt bucket objects for post-transcode"
  deletion_window_in_days = 10
}

resource "aws_s3_bucket_server_side_encryption_configuration" "post_example" {
  bucket = aws_s3_bucket.myp_post_transcode.id

  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = aws_kms_key.s3key_post.arn
      sse_algorithm     = "aws:kms"
    }
  }
}

# Versioning for pre-transcode bucket
resource "aws_s3_bucket_versioning" "versioning_pre_example" {
  bucket = aws_s3_bucket.myp_pre_transcode.id
  versioning_configuration {
    status = "Enabled"
  }
}

# Versioning for post-transcode bucket
resource "aws_s3_bucket_versioning" "versioning_post_example" {
  bucket = aws_s3_bucket.myp_post_transcode.id
  versioning_configuration {
    status = "Enabled"
  }
}
