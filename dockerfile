FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /restapi

FROM alpine:latest

WORKDIR /

COPY --from=builder /restapi /restapi

EXPOSE 8080

CMD ["/restapi"]