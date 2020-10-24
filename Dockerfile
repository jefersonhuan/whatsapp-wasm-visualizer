FROM golang:1.15 AS builder

WORKDIR /usr/go/whatsapp-visualizer/serve
COPY ./serve .

ENV GOOS=linux
ENV CGO_ENABLED=0
RUN go build start.go

WORKDIR /usr/go/whatsapp-visualizer
COPY . .

ENV GOOS=js
ENV GOARCH=wasm
RUN go build -o main.wasm

FROM alpine:3.9

WORKDIR /var/app

COPY --from=builder /usr/go/whatsapp-visualizer/serve/ .
COPY --from=builder /usr/go/whatsapp-visualizer/main.wasm .

CMD ["/var/app/start"]