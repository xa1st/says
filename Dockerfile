# 使用golang:alpine镜像作为构建环境
FROM golang:alpine AS build

# 设置工作目录为/app
WORKDIR /app
# 将当前目录下的所有文件复制到工作目录中
COPY . .
# 更新go模块并进行构建，-ldflags参数用于减小二进制文件大小
RUN go mod tidy && go build -ldflags="-s -w" -o says .

# 切换到alpine:latest镜像作为最终的运行环境
FROM alpine:latest

# 设置工作目录为/app
WORKDIR /app
# 从构建阶段的镜像中复制编译好的`says`二进制文件到当前目录
COPY --from=build /app/says .

# 暴露8080端口，供外部访问
EXPOSE 3000

# 指定容器启动时执行的命令
CMD ["./says"]