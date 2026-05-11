## Build
FROM golang:1.25 AS build
WORKDIR /app

# Copy the Go Modules manifests
# cache deps before building and copying source so that we don't need to re-download
# as much and so that source changes don't invalidate our downloaded layer
COPY go.mod go.sum ./

RUN go mod download 
ADD . /app
RUN env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o todo-list ./cmd/main.go 


## Deploy
FROM alpine
USER root

# set correct timezone
RUN apk --no-cache add tzdata && rm -rf /var/cache/apk/*
ENV TZ=Europe/Moscow

WORKDIR /app
COPY --from=build /app/todo-list todo-list
COPY --from=build /app/web web
COPY --from=build /app/migrations migrations
COPY --from=build /app/parameters.yaml parameters.yaml
ENTRYPOINT ["./todo-list"]
