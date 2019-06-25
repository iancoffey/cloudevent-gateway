FROM golang:alpine
ENV GOBIN /go/bin

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
RUN go get -u github.com/golang/dep/...
RUN dep ensure
RUN go build -o main ./cmd/gateway/...
RUN adduser -S -D -H -h /app appuser
USER appuser
CMD ["./main"]


