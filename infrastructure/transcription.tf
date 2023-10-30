resource "aws_iam_policy" "transcribe_policy" {
  name        = "TranscribeStartJobPolicy"
  description = "Policy to allow user to start Transcribe jobs"
  policy      = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
              "transcribe:StartTranscriptionJob",
              "transcribe:GetTranscriptionJob"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_user_policy_attachment" "attach_transcribe_policy" {
  user       = "go-audibly"  // Replace with the actual IAM username
  policy_arn = aws_iam_policy.transcribe_policy.arn
}
