# 🎙️ File Watcher to AWS Transcriber

![Build](https://github.com/kmesiab/go-audibly/actions/workflows/go.yml/badge.svg)

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

### Clone the Repository

```bash
git clone https://github.com/kmesiab/go-audibly.git
```

### Initialize Infrastructure

Navigate to the `./infrastructure` folder and run:

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
./my_app
```

Use the `Makefile` for common tasks:

- `make init` : Initialize Terraform in `./infrastructure`
- `make plan` : Plan Terraform changes in `./infrastructure`
- `make apply`: Apply Terraform changes in `./infrastructure`
- `make build`: Build the Go application
- `make test` : Run Go tests
- `make fumpt`: Run fumpt to format code
- `make lint` : Run golangci-lint
