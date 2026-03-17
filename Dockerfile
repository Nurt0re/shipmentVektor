FROM golang:1.25.4-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /build/app .
COPY --from=builder /build/internal/migrations ./internal/migrations
COPY --from=builder /build/.env .

EXPOSE 8080

CMD ["./app"]
