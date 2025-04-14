FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY .env .                

RUN go build -o user-service-app ./cmd/user/main.go

EXPOSE 6060

CMD ["./user-service-app"]

