FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o status-checker

FROM alpine:latest
RUN apk add --no-cache curl busybox
WORKDIR /app
COPY --from=builder /app/status-checker /app/status-checker

CMD [ "/app/status-checker" ]

