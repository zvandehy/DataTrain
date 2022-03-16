
#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
COPY ./app.env .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v server.go

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
COPY ./app.env .
ENTRYPOINT /app
LABEL Name=datatrain-graphql Version=1.0.0
EXPOSE 8080