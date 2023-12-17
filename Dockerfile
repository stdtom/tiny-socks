FROM golang:1.20 AS build
ENV GOFLAGS="-mod=vendor"

COPY . /go/src/github.com/stdtom/tiny-socks

WORKDIR /go/src/github.com/stdtom/tiny-socks/cmd/tiny-socks
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/tiny-socks .


FROM scratch AS runtime
USER 1234
WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/tiny-socks /

CMD ["/tiny-socks"]
