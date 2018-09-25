FROM golang:alpine as builder
ADD . /go/src/github.com/VMitov/payments
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static" -w -s' -o /go/bin/app github.com/VMitov/payments/cmd/payments

FROM scratch
COPY --from=builder /go/bin/app /
ENTRYPOINT ["/app"]
