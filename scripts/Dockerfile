FROM golang:1.23-alpine3.21 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY keys keys
RUN mkdir cmd
RUN mkdir pkg
COPY cmd/api cmd/api
COPY pkg/model pkg/model
COPY pkg/logger pkg/logger
COPY pkg/storage pkg/storage
COPY pkg/service pkg/service
COPY pkg/server pkg/server
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api/main.go 


FROM alpine:3.21
WORKDIR /app
COPY --from=build /app/api api
CMD ["api", "-c", "config.json"]