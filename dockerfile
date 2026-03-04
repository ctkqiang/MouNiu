# 使用官方Go镜像作为构建环境
FROM golang:1.25-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o mouniu main.go

# 使用轻量级的Alpine镜像作为运行环境
FROM alpine:3.18

# 设置工作目录
WORKDIR /app

# 复制构建好的应用
COPY --from=builder /app/mouniu .

# 复制配置文件
COPY internal/config/ ./internal/config/

# 暴露端口
EXPOSE 8000

# 设置环境变量默认值
ENV APP_NAME=MouNiu
ENV QUESTDB_DB=qdb
ENV QUESTDB_HOST=localhost
ENV QUESTDB_USER=admin
ENV QUESTDB_PASS=quest
ENV QUESTDB_PORT=8812

# 运行应用
CMD ["./mouniu"]