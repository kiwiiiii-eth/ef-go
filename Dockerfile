# 使用官方 Golang 映像作為構建階段
FROM golang:1.21-alpine AS builder

# 設置工作目錄
WORKDIR /app

# 安裝 git（某些依賴可能需要）
RUN apk add --no-cache git

# 複製 go mod 文件
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製源代碼
COPY . .

# 構建應用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# 使用輕量級映像運行
FROM alpine:latest

# 安裝 CA 證書（HTTPS 請求需要）
RUN apk --no-cache add ca-certificates tzdata

# 設置時區為台北
ENV TZ=Asia/Taipei

# 創建應用目錄
WORKDIR /root/

# 從構建階段複製二進制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 運行應用
CMD ["./main"]
