FROM golang:1.16 AS build
ENV CGO_ENABLED=0
WORKDIR /go/src

COPY go.mod .
COPY go.sum .
COPY main.go .
COPY pkg ./pkg

RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o openapi .

FROM scratch AS runtime
ENV GIN_MODE=release
EXPOSE 3000/tcp
ENTRYPOINT ["./openapi"]

COPY --from=build /pkg/core/src/openapi ./
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
