# Makefile for running Terraform and Go commands

# 🌍 Run app
run:
	@echo "Starting!"
	source .env && go build . && go run .

# 🌍 Terraform targets
init:
	@echo "🌱 Initializing Terraform in /infrastructure..."
	source .env && cd ./infrastructure && terraform init

deploy:
	@echo "💣 Deploying infrastructure."
	source .env && cd ./infrastructure && terraform init && terraform plan -out=tfplan && terraform apply -auto-approve tfplan

destroy:
	@echo "💣 Destroying Terraform resources in /infrastructure..."
	source .env && cd ./infrastructure && terraform destroy

# 🏗 Go build and test targets
build:
	@echo "🛠 Building Go project..."
	go build -o .

test:
	@echo "🚀 Running Go tests..."
	go test ./...

# 🌈 All-in-one linting
lint:
	@echo "🔍 Running all linters..."
	golangci-lint run && markdownlint README.md

# 🌈 All-in-one build, test, and lint
all: build test lint
	@echo "🎉 Done!"

.PHONY: build test lint all
