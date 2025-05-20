FROM golang:1.24

WORKDIR /app
RUN go version
ENV $GOPATH=/

COPY . .

RUN go mod download
RUN go build -o nameservice ./cmd/main.go

EXPOSE 8080

CMD ["./nameservice"]