# 评测用镜像：构建并直接启动项目真实入口。
FROM golang:1.22
WORKDIR /app
COPY go.mod ./
COPY . .
RUN go build -o /usr/local/bin/app ./cmd/forex
CMD ["/usr/local/bin/app"]

# 多架构交叉构建示例（如需交付双架构镜像）：
# docker buildx build --platform linux/arm64,linux/amd64 -f benzhi.Dockerfile -t <image> .
