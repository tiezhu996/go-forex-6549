# Evaluation image for the repository's real one-shot CLI entrypoint.
FROM golang:1.22
WORKDIR /app
COPY go.mod ./
COPY . .
RUN go build -o /usr/local/bin/forex ./cmd/forex
CMD ["/usr/local/bin/forex"]

# Multi-architecture delivery example:
# docker buildx build --platform linux/arm64,linux/amd64 -f benzhi.Dockerfile -t <image> .
