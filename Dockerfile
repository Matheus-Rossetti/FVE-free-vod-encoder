FROM golang:1.26.1 AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o frevod ./src/

FROM alpine:latest

RUN apk add --no-cache \
    ffmpeg \
    ca-certificates

WORKDIR /app

COPY --from=builder /app/frevod .

ENTRYPOINT ["./frevod"]