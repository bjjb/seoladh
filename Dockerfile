FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/seoladh
COPY . .
RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/seoladh
FROM scratch
COPY --from=builder /go/bin/seoladh /bin/seoladh
EXPOSE 12345
ENTRYPOINT ["/bin/seoladh"]
