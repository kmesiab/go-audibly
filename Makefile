# Makefile for running Terraform and Go commands
# 🌍 Run app
run:
	@echo "Starting!"
	go build . && go run .

# 🌍 Terraform targets
init:
	@echo "🌱 Initializing Terraform in /infrastructure..."
	cd ./infrastructure && terraform init

plan:
	@echo "🔍 Planning Terraform changes in /infrastructure..."
	cd ./infrastructure && terraform plan

apply:
	@echo "✅ Applying Terraform changes in /infrastructure..."
	cd ./infrastructure && terraform apply

destroy:
	@echo "💣 Destroying Terraform resources in /infrastructure..."
	cd ./infrastructure && terraform destroy

deploy:
	@echo "💣 Deploying infrastructure."
	cd ./infrastructure && terraform plan && terraform deploy

# 🏗 Go build and test targets
build:
	@echo "🛠 Building Go project..."
	go build -o .

test:
	@echo "🚀 Running Go tests..."
	go test ./...

# 🔍 Linting with golangci-lint
lint:
	@echo "🔍 Running golangci-lint..."
	golangci-lint run

# 📜 Lint README
readme-lint:
	@echo "📝 Linting README..."
	markdownlint README.md

# 🌈 All-in-one
all: build test lint readme-lint
	@echo "🎉 Done!"

.PHONY: build test lint readme-lint all
