# build stage
FROM golang:latest AS builder
LABEL maintainer = "hilli@github.com"
RUN mkdir -p /go/src/finance-influxdb
WORKDIR /go/src/finance-influxdb
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /finance-influxdb .

# final stage
FROM alpine:latest
COPY --from=builder /finance-influxdb ./
RUN chmod +x ./finance-influxdb
ENTRYPOINT ["./finance-influxdb"]
