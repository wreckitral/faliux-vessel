FROM golang:1.25-alpine AS builder

WORKDIR /workspace

COPY . .

WORKDIR /workspace/faliux-vessel-cli-consignment

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o faliux-vessel-cli-consignment

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /workspace/faliux-vessel-cli-consignment/consignment.json ./
COPY --from=builder /workspace/faliux-vessel-cli-consignment/faliux-vessel-cli-consignment .

ENTRYPOINT ["./faliux-vessel-cli-consignment"]
