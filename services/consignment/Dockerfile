FROM golang:1.25-alpine AS builder

WORKDIR /workspace

COPY go.work go.work.sum ./

COPY services/consignment/go.mod services/consignment/go.sum ./services/consignment/
COPY services/vessel/go.mod services/vessel/go.sum ./services/vessel/
COPY clients/cli-consignment/go.mod clients/cli-consignment/go.sum ./clients/cli-consignment/

WORKDIR /workspace/services/consignment
RUN go mod download

WORKDIR /workspace
COPY services/consignment/ ./services/consignment/
COPY services/vessel/generated/ ./services/vessel/generated/

WORKDIR /workspace/services/consignment
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o consignment-service .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /workspace/services/consignment/consignment-service /app/consignment-service
ENTRYPOINT ["/app/consignment-service"]
