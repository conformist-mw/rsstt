FROM golang:1.21-bookworm AS stage1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o rsstt

from scratch AS export-stage
COPY --from=stage1 /app/rsstt .
