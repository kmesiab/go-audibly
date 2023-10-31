# 🎙️ File Watcher to AWS Transcriber

![Terraform](https://img.shields.io/badge/terraform-%235835CC.svg?style=for-the-badge&logo=terraform&logoColor=white)
![AWS](https://img.shields.io/badge/AWS-%23FF9900.svg?style=for-the-badge&logo=amazon-aws&logoColor=white)

The Terraform in the `./infrastructure` folder automates the setup of AWS resources for handling pre-transcode and post-transcode 
S3 buckets, AWS Transcribe access, and AWS KMS encryption. It also sets up an IAM user with the necessary permissions to interact 
with these resources.

## Components :building_construction:

- **S3 Buckets** :file_folder:: Two S3 buckets are created, `myp-pre-transcode` for pre-transcode and `myp-post-transcode` for post-transcode.
- **Bucket Policies** :lock:: Grants AWS Transcribe service the permission to `GetObject` from both buckets.
- **AWS KMS Encryption** :key:: Keys are generated for each bucket to enable server-side encryption.
- **Versioning** :clock3:: Enables versioning for both buckets.
- **IAM User** :bust_in_silhouette:: Creates an IAM user `go-audibly` with S3 full access and KMS data key generation permissions.
- **Transcribe Policy** :microphone:: Allows the user to start and get Transcribe jobs.

## Requirements :clipboard:

1. Terraform installed
2. AWS CLI configured with appropriate access

## Usage :rocket:

1. Clone the repository containing the Terraform script.
2. Navigate to the directory and run `terraform init` to initialize the Terraform configuration.
3. Run `terraform apply` to create the AWS resources. Confirm by typing `yes` when prompted.

## Customization :wrench:

- Update the `bucket` names in `aws_s3_bucket` resources if you want different bucket names.
- Update the `tags` to match your environment labeling.
- Update the `policy` under `aws_iam_policy` if you need to fine-tune permissions.

Run `terraform apply` again after making changes.

## Outputs :outbox_tray:

- `aws_access_key_id`: The IAM user's access key ID.
- `aws_secret_access_key`: The IAM user's secret access key.
- `secret_file_path`: Path to the file containing the secret access key.

:warning: **Note**: The secret access key is sensitive information. Handle with care.

## Cleanup :broom:

Run `terraform destroy` to remove all resources created by this script. Confirm by typing `yes` when prompted.
