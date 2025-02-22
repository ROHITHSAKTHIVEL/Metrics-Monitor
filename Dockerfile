FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN apk add --no-cache postgresql-client

COPY . .

RUN go build -o metrics-monitor main.go

EXPOSE 8888

CMD ["./metrics-monitor"]
