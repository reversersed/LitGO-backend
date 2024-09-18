FROM golang:alpine AS builder

WORKDIR /usr/local/go/src/

ADD . .

RUN export GODEBUG=http2=0
RUN go clean --modcache
RUN go build -mod=readonly -o app cmd/genre/main.go

FROM alpine:latest

COPY --from=builder /usr/local/go/src/app /

CMD ["/app"]