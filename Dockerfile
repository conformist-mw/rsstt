FROM golang:1.21-alpine AS build
RUN apk --no-cache add build-base gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64v8 go build -o rsstt

FROM alpine:latest

COPY --from=build /app/rsstt /

ENTRYPOINT [ "/rsstt" ]
