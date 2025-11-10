FROM golang:1.25-alpine AS builder

WORKDIR /workspace/services/vessel

COPY services/vessel/go.mod services/vessel/go.sum ./
RUN go mod download

COPY services/vessel/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o vessel-service .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /workspace/services/vessel/vessel-service .

ENTRYPOINT ["./vessel-service"]
