provider "aws" {
  region = "us-west-2"
}

resource "aws_s3_bucket" "myp_pre_transcode" {
  bucket = "myp-pre-transcode"

  tags = {
    Name        = "myp"
    Environment = "dev"
  }
}
