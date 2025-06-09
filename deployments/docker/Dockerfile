FROM golang:1.20-alpine AS builder

WORKDIR /app

# 设置 Go 模块代理
ENV GOPROXY=https://goproxy.cn,direct

# 复制源代码
COPY . .

# 构建应用
RUN go build -o eden-ops cmd/eden-ops/main.go

FROM alpine:latest

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/eden-ops .
COPY configs/config.yaml.example /etc/eden-ops/config.yaml

# 设置时区
RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

EXPOSE 8080

ENTRYPOINT ["./eden-ops"] 