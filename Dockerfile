FROM golang:1.17-buster AS builder

ENV GOPATH=/

WORKDIR /app/
COPY . /app/

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o /apiserver ./cmd/apiserver

FROM alpine:latest

# copy compiled go app
COPY --from=builder /apiserver /apiserver

CMD ["./apiserver"]