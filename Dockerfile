ARG GOLANG_VERSION=1.17
FROM golang:${GOLANG_VERSION}-buster as builder
ARG GOPROXY=https://goproxy.cn
RUN sed -i 's/ports.ubuntu.com/mirror.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list
WORKDIR ${GOPATH}/src/github.com/projectxpolaris/youphoto

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ${GOPATH}/bin/youphoto ./main.go

FROM debian:buster-slim

COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /etc/ssl/certs /etc/ssl/certs


COPY --from=builder /go/bin/youphoto /usr/local/bin/youphoto



ENTRYPOINT ["/usr/local/bin/youphoto","run"]

