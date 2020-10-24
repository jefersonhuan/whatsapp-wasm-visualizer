FROM golang:1.15 AS builder

WORKDIR /usr/go/whatsapp-visualizer/server
COPY ./server .

ENV GOOS=linux
ENV CGO_ENABLED=0
RUN go build serve.go

WORKDIR /usr/go/whatsapp-visualizer
COPY . .

ENV GOOS=js
ENV GOARCH=wasm
RUN go build -o main.wasm

FROM alpine:3.9

WORKDIR /var/app

COPY --from=builder /usr/go/whatsapp-visualizer/server/ .
COPY --from=builder /usr/go/whatsapp-visualizer/main.wasm .

RUN ls -l

CMD ["/var/app/serve"]