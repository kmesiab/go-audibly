# Existing S3 bucket for pre-transcribe
resource "aws_s3_bucket" "audio_bucket" {
  bucket = "myp-audio-bucket"

  lifecycle {
    prevent_destroy = false
  }

  tags = {
    Name        = "myp"
    Environment = "dev"
  }
}

# Existing S3 bucket for post-transcribe
resource "aws_s3_bucket" "transcript_bucket" {
  bucket = "myp-transcript-bucket"

  lifecycle {
    prevent_destroy = false
    ignore_changes = [policy]
  }

  tags = {
    Name        = "myp"
    Environment = "dev"
  }
}

# Policy to allow Transcribe to read from the pre-transcribe bucket
resource "aws_s3_bucket_policy" "audio_bucket_policy" {
  bucket = aws_s3_bucket.audio_bucket.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect    = "Allow",
        Principal = { Service = "transcribe.amazonaws.com" },
        Action    = "s3:*",
        Resource  = "${aws_s3_bucket.audio_bucket.arn}/*"
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
  bucket = aws_s3_bucket.audio_bucket.id

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
  bucket = aws_s3_bucket.transcript_bucket.id

  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = aws_kms_key.s3key_post.arn
      sse_algorithm     = "aws:kms"
    }
  }
}

# Versioning for pre-transcribe bucket
resource "aws_s3_bucket_versioning" "versioning_pre_bucket" {
  bucket = aws_s3_bucket.audio_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}

# Versioning for post-transcribe bucket
resource "aws_s3_bucket_versioning" "versioning_post_bucket" {
  bucket = aws_s3_bucket.transcript_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}
