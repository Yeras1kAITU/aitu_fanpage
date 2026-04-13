FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/server

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata wget

WORKDIR /app

COPY --from=builder /app/main .

RUN mkdir -p /data/uploads /data/temp \
    /data/uploads/images /data/uploads/videos \
    /data/uploads/documents /data/uploads/thumbnails

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]