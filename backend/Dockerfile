FROM golang:1.24.12-alpine3.23 AS builder
RUN mkdir -p /var/go
WORKDIR /var/go
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

FROM alpine:3.23.3
WORKDIR /var/go
COPY --from=builder /var/go .
EXPOSE 8080
ENTRYPOINT [ "./main" ]