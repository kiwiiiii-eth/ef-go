.PHONY: help build run test clean docker-build docker-run

help: ## 顯示幫助信息
	@echo "可用的命令："
	@echo "  make build        - 構建應用程式"
	@echo "  make run          - 運行應用程式"
	@echo "  make test         - 運行測試"
	@echo "  make clean        - 清理構建文件"
	@echo "  make docker-build - 構建 Docker 映像"
	@echo "  make docker-run   - 運行 Docker 容器"

build: ## 構建應用程式
	go build -o bin/vpp-api ./cmd/api

run: ## 運行應用程式
	go run ./cmd/api/main.go

test: ## 運行測試
	go test -v ./...

clean: ## 清理構建文件
	rm -rf bin/
	go clean

docker-build: ## 構建 Docker 映像
	docker build -t vpp-go:latest .

docker-run: ## 運行 Docker 容器
	docker run -p 8080:8080 --env-file .env vpp-go:latest

deps: ## 下載依賴
	go mod download
	go mod tidy

lint: ## 運行代碼檢查
	golangci-lint run

format: ## 格式化代碼
	go fmt ./...
	goimports -w .
