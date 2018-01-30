# ApiMock
## Install
```
go get -u github.com/gorilla/mux
```
```
go get -u github.com/google/uuid
```
## Build 
```
go build
```
## Test
Request :
```
wget http://localhost/mock
```
Response :
```
{
  "uuid": "a7b2b2c9-b06a-458e-8467-4f56e000cfa2",
  "message": "Hello world !",
  "headers": [
    {
      "value": "Accept:text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
    },
    {
      "value": "Accept-Encoding:gzip, deflate"
    },
    {
      "value": "Upgrade-Insecure-Requests:1"
    },
    {
      "value": "Cache-Control:max-age=0"
    },
    {
      "value": "User-Agent:Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:57.0) Gecko/20100101 Firefox/57.0"
    },
    {
      "value": "Accept-Language:fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3"
    },
    {
      "value": "Dnt:1"
    },
    {
      "value": "Connection:keep-alive"
    }
  ]
}
```
