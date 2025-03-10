# development
FROM golang:1.24 AS development

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8081

CMD ["go", "run", "main.go"]

# build
FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8081

CMD ["./main"]
