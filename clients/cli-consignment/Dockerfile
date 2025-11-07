FROM golang:1.25-alpine AS builder

WORKDIR /workspace/clients/cli-consignment

COPY clients/cli-consignment/go.mod clients/cli-consignment/go.sum ./
RUN go mod download

COPY clients/cli-consignment/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o cli-consignment .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /workspace/clients/cli-consignment/consignment.json ./
COPY --from=builder /workspace/clients/cli-consignment/cli-consignment .

ENTRYPOINT ["./cli-consignment"]
