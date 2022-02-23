#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
COPY ./config.yaml .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v server.go

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
COPY ./config.yaml .
ENTRYPOINT /app
LABEL Name=faceapi Version=1.0.3
EXPOSE 8080
