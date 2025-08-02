FROM golang:1.24.1-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/app ./cmd/main.go

FROM alpine:latest
COPY --from=builder /build/app /application/app
EXPOSE 8080
ENTRYPOINT ["/application/app"]
