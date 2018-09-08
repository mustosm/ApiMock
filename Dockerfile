FROM golang:1.10.3 as builder
WORKDIR /home/must/work/go/src/ApiMock
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /home/must/work/go/src/ApiMock/ApiMock .
COPY --from=builder /home/must/work/go/src/ApiMock/ApiMock.properties .
CMD ["./ApiMock"]