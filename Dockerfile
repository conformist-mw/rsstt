ARG TARGETOS
ARG TARGETARCH

FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o rsstt

FROM alpine:latest

COPY --from=builder /app/rsstt /

ENTRYPOINT [ "/rsstt" ]
