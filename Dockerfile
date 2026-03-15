FROM golang:1.25.7-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GCO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w" \
    -trimpath \
    -o user-api .

FROM alpine:3.23.3

LABEL org.opencontainers.image.source=https://github.com/ishs-cloud-computing/user-api

COPY --from=builder /app/user-api /user-api

EXPOSE 8080

ENTRYPOINT [ "./user-api" ]
