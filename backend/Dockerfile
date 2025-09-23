FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates && update-ca-certificates
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN if [ -f go.work ]; then go work sync; fi

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/app -trimpath -tags netgo -ldflags '-s -w' ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata && update-ca-certificates

WORKDIR /app
COPY --from=builder /out/app ./main

EXPOSE 8080
CMD ["./main"]