FROM golang:1.21.4

WORKDIR /cmd/bookmarket

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v cmd/bookmarket/main.go

EXPOSE 8080

CMD ["./main"]