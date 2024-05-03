FROM golang:1.22 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o parser cmd/parser/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/parser /app/parser
COPY config/docker.yaml /app/docker.yaml
WORKDIR /app
CMD ["./parser", "-config=docker.yaml"]