# 第一阶段：构建 Go 程序
FROM golang:1.20-alpine AS builder

# 安装必要的构建依赖
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 将当前目录的所有文件复制到工作目录
COPY . .

# 下载并安装依赖
RUN go mod tidy

# 编译程序
RUN go build -o main .

# 第二阶段：创建一个最小的运行时环境
FROM alpine:latest

# 设置环境变量
ENV LLM_API_URL=https://one-api.skybyte.me/v1/chat/completions
ENV LLM_MODEL_NAME=gemini-1.5-flash
ENV LLM_API_KEY=sk-rLEbROmxDX3MfRq2155e4138BbAb48699e8c6203D492FaE1

# 设置工作目录
WORKDIR /root/

# 从 builder 镜像复制编译好的二进制文件到当前工作目录
COPY --from=builder /app/main .

# 运行程序
CMD ["./main"]

