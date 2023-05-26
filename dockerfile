FROM golang:alpine as builder
WORKDIR /go/delivery
COPY . .
RUN go build -o ably_log_recorder .

FROM alpine
COPY --from=builder /go/delivery/ably_log_recorder .
CMD ["./ably_log_recorder"]