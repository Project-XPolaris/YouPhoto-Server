# 构建阶段
FROM golang:1.24-alpine AS builder

# 添加元数据
LABEL maintainer="ProjectXPolaris"
LABEL description="YouPhoto Server - 图片管理服务"

# 设置构建参数
ARG GOPROXY=https://goproxy.cn
ENV GOPROXY=${GOPROXY}
ENV CGO_ENABLED=0
ENV GOOS=linux

# 设置工作目录
WORKDIR /build

# 首先复制依赖文件以利用缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o youphoto ./main.go

# 运行阶段
FROM alpine:latest

# 安装基础依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN adduser -D youphoto
USER youphoto

# 复制构建产物
COPY --from=builder /build/youphoto /usr/local/bin/

# 设置工作目录
WORKDIR /app

# 暴露端口（如果需要的话，请根据实际情况修改）
EXPOSE 8000

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8000/health || exit 1

# 运行应用
ENTRYPOINT ["/usr/local/bin/youphoto", "run"]

