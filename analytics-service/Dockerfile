FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN go build -o /bin/activity-analytics-service ./cmd/app

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /bin/activity-analytics-service /bin/activity-analytics-service

EXPOSE 8084
CMD ["/bin/activity-analytics-service"]
