# gauche

[![tjvc](https://circleci.com/gh/tjvc/gauche.svg?style=shield)](https://circleci.com/gh/tjvc/gauche)
[![Go Report Card](https://goreportcard.com/badge/github.com/tjvc/gauche)](https://goreportcard.com/report/github.com/tjvc/gauche)

> gauche
>
> 1 b : crudely made or done

Gauche ("Go cache") is a simple HTTP in-memory key-value store.

Roadmap: https://www.pivotaltracker.com/n/projects/2487327

## Usage

### Build from source

```
go mod download
go build
./gauche
```

### Docker

```
docker build -t gauche .
docker run -p 8080:8080 gauche
```

### PUT key

```bash
$ curl -i -X PUT -d 'value' http://localhost:8080/key
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Sun, 14 Feb 2021 10:19:36 GMT
Content-Length: 5

value
```

### GET key

```bash
$ curl -i http://localhost:8080/key
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Sun, 14 Feb 2021 10:20:31 GMT
Content-Length: 5

value
```

### DELETE key

```bash
$ curl -i -X DELETE http://localhost:8080/key
HTTP/1.1 204 No Content
Date: Sun, 14 Feb 2021 10:22:22 GMT
```

### GET keys

```bash
$ curl -i http://localhost:8080
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Sun, 14 Feb 2021 10:23:11 GMT
Content-Length: 24

key1
key2
key3
```