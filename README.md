
# Marvel Characters Golang REST Api

A RESTful api that's based on [Marvel API](https://developer.marvel.com/) that returns:
- all the Marvel characters ids in a JSON array of numbers.
- the id, name and description of a character.

## Requirements
- Sign up for a free **Marvel Developer API** at https://developer.marvel.com/.
- [Go version 1.15](https://golang.org/dl/) or greater.
- Install [mux router](https://github.com/gorilla/mux).

## Quick Start

``` bash
# Install mux router
go get -u github.com/gorilla/mux
```

``` bash
# Build and run the rest api
cd "$GOPATH"/src/github.com/albertleng/restapi || exit
go build -o restapi
echo "Go build done"
./restapi
```

## Endpoints

### Get all Marvel character ids
``` bash
GET /characters
```

### Get id, name and description of a character
``` bash
GET /characters/{characterId}
```


