FROM golang:latest as BUILDER

WORKDIR /go/src/

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/main .

FROM alpine:latest

WORKDIR /go/bin

COPY --from=builder /go/bin/main .

CMD ["/go/bin/main"]