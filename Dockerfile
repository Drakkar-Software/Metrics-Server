FROM golang:1.13 as builder

ADD . /${GOPATH}/src/github.com/Drakkar-Software/Metrics-Server
WORKDIR /${GOPATH}/src/github.com/Drakkar-Software/Metrics-Server

RUN go get -u github.com/golang/dep/cmd/dep \
    && dep ensure \
    && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o server .

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/Drakkar-Software/Metrics-Server/server /server/
WORKDIR /server

EXPOSE 8080

# Start up
CMD ["./server"]