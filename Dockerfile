FROM golang:1.10

WORKDIR /home/must/go/src/ApiMock
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]