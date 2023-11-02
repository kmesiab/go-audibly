# 🎙️ File Watcher to AWS Transcriber

![Build](https://github.com/kmesiab/go-audibly/actions/workflows/go.yml/badge.svg)

![Golang](https://img.shields.io/badge/Go-00add8.svg?labelColor=171e21&style=for-the-badge&logo=go)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Terraform](https://img.shields.io/badge/terraform-%235835CC.svg?style=for-the-badge&logo=terraform&logoColor=white)
![AWS](https://img.shields.io/badge/AWS-%23FF9900.svg?style=for-the-badge&logo=amazon-aws&logoColor=white)

![Logo](./assets/logo-go-aud-small.png)

## 📋 Table of Contents

- [Overview](https://github.com/kmesiab/go-audibly/README.md#overview)
- [Prerequisites](https://github.com/kmesiab/go-audibly/README.md#prerequisites)
- [Setup](https://github.com/kmesiab/go-audibly/README.md#setup)
- [Usage](https://github.com/kmesiab/go-audibly/README.md#usage)

## 📜 Overview

This application watches a folder for new files. Upon detecting one, it uploads
the file to an AWS S3 bucket and starts an AWS Transcribe job. The transcribed
results are saved in another desktop folder.

## 🛠️ Prerequisites

- Go installed
- AWS CLI installed and configured
- Terraform installed

## 🚀 Setup

### 🌱 Environment Variables (.env file)

#### Description

The .env file contains environment variables that are crucial for
the application's configuration. These variables are used to specify
various settings like AWS configurations, folder paths for audio and
transcript files, etc.

#### Variables

- APP_NAME: The name of the application.
- ENV: Specifies the environment in which the application is running.
Set this according to your SDLC.
- WATCH_FOLDER: The folder path to watch for new audio files.
- TRANSCRIPT_FOLDER: The folder path where transcribed text files will
be saved.
- PROCESSED_AUDIO_FOLDER: The folder path where processed audio files
will be saved.
- AWS_S3_INPUT_BUCKET_NAME: AWS S3 bucket name for input (audio) files.
- AWS_S3_OUTPUT_BUCKET_NAME: AWS S3 bucket name for output (transcripts) files.
- AWS_REGION: AWS region.
- AWS_ACCESS_KEY_ID: Your AWS access key.
- AWS_SECRET_ACCESS_KEY: Your AWS secret access key.

```bash
export APP_NAME="go-audibly"
export ENV="local"

export WATCH_FOLDER="./pre-processed"
export PROCESSED_AUDIO_FOLDER="./post-processed"
export TRANSCRIPT_FOLDER="./transcripts"

export AWS_S3_INPUT_BUCKET_NAME="audio"
export AWS_S3_OUTPUT_BUCKET_NAME="transcripts"
export AWS_REGION="us-west-2"

export AWS_ACCESS_KEY_ID="<SECRET>"
export AWS_SECRET_ACCESS_KEY="<SECRET>"
```

#### Usage

Create a .env file at the root of your project and copy the above variables
into it. Make sure to replace the `<SECRET>` placeholders with your actual
AWS credentials.

⚠️ *Security Note: Never commit your .env file containing sensitive AWS
credentials into version control. Add .env to your .gitignore file.*

### Clone the Repository

```bash
git clone https://github.com/kmesiab/go-audibly.git
```

### Initialize Infrastructure

From the root folder, run:

```bash
make init
make plan
make apply
```

### Build and Test Application

In the project root, run:

```bash
make build
make test
```

## 🎯 Usage

Run the application:

```bash
./go-audibly
```

Use the `Makefile` for common tasks:

- `make init` : Initialize Terraform in `./infrastructure`
- `make plan` : Plan Terraform changes in `./infrastructure`
- `make apply`: Apply Terraform changes in `./infrastructure`
- `make build`: Build the Go application
- `make test` : Run Go tests
- `make fumpt`: Run fumpt to format code
- `make lint` : Run golangci-lint

## Terraform 📘

🔗 [Terraform Instructions](https://github.com/kmesiab/go-audibly/blob/main/infrastructure/README.md)
