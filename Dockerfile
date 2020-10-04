FROM golang:1.13.10 as builder
COPY ./src /go/src/app
WORKDIR /go/src/app
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o monitor .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/app/monitor .
COPY --from=builder /go/src/app/config.ini .
CMD [ "/root/monitor" ]