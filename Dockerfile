FROM golang:1.16-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o scuba-divers-app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/scuba-divers-app .

EXPOSE 8080

CMD ["./scuba-divers-app"]
