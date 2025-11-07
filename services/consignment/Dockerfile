FROM golang:1.25-alpine AS builder

WORKDIR /workspace/services/consignment

COPY services/consignment/go.mod services/consignment/go.sum ./
RUN go mod download

COPY services/consignment/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o consignment-service .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /workspace/services/consignment/consignment-service .

ENTRYPOINT ["./consignment-service"]
