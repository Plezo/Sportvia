FROM golang:1.22.1-bookworm

WORKDIR /sportvia
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o sportvia ./cmd/app

EXPOSE 8080

CMD ./sportvia -bind=0.0.0.0:8080 -env=production
