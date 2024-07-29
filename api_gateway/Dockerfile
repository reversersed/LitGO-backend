FROM golang:alpine AS builder

WORKDIR /usr/local/go/src/

ADD . .

RUN go clean --modcache
RUN go build -mod=readonly -o app cmd/gateway/main.go

FROM alpine:latest

COPY --from=builder /usr/local/go/src/app /

EXPOSE 9000

CMD ["/app"]