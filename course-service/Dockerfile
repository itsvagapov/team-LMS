FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /course-service \
    ./cmd/app

FROM alpine:3.23

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S app -G app

WORKDIR /app

COPY --from=builder /course-service ./course-service

USER app

EXPOSE 8081

ENTRYPOINT ["./course-service"]
