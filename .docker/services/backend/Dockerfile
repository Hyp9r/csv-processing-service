FROM golang:1.23.3-alpine

WORKDIR /app

COPY ./go.mod .

RUN go mod download

COPY ./ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app cmd/webserver/main.go

EXPOSE 8081

ENTRYPOINT ["/go/bin/app"]