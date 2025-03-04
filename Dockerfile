# Build stage
FROM golang:1.24.0-alpine3.20 AS builder
WORKDIR /app
COPY . .
ENV GOPROXY=https://goproxy.cn,direct

RUN go build -o main *.go

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY start.sh .
COPY wait-for.sh .

EXPOSE 8888
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]