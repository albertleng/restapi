# Marvel Characters Golang REST Api

A RESTful api that's based on [Marvel API](https://developer.marvel.com/) that returns:

- all the Marvel character `ids` in a JSON array of numbers.
- the `id`, `name` and `description` of a character.

## Requirements

- Sign up for a free **Marvel Developer API** at https://developer.marvel.com/.
- [Go version 1.15](https://golang.org/dl/) or greater.
- Install [mux router](https://github.com/gorilla/mux) to implement a request router and dispatcher for matching
  incoming requests to their respective handler.
- Install [Go CORS handler](https://github.com/rs/cors) to handle cors requests.
- This api is developed and verified in [Ubuntu 18.04.5 LTS](https://releases.ubuntu.com/18.04/)

## Quick Start

### Install this Marvel REST Api

``` bash
go get -u github.com/albertleng/restapi
```

### Install mux router

``` bash
go get -u github.com/gorilla/mux
```

### Install CORS handler

``` bash
go get -u github.com/rs/cors
```

### Environment Variables

To run this project, you will need to add the following environment variables:

`MARVEL_API_PRIVATE_KEY`

`MARVEL_API_PUBLIC_KEY`

### Build and run the rest api

Run the `build.sh` to build and run the rest api

``` bash
./build.sh
```

or

``` bash
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

## Tests

`main_test.go` contains tests for `/characters` and `/characters/{characterId}`.

``` bash
# run test
go test -v
```

## Caching Strategy

This api uses a simple text file `ids.txt` to store the `character ids`.

In response to a `/characters` call, there will be two scenarios:

### 1. `ids.txt` exists

The api reads all the `character ids` from `ids.txt` into an integer array, `ids`.

- If it is the first call after the api starts, the api makes calls to `https://gateway.marvel.com/v1/public/characters`
  with `offset` starting with length of `ids` (and increment of 100 in subsequent calls) and `limit` of 100 until there
  is no character returned. These character ids are appended to both `ids.txt` and `ids`. The updated `ids` is returned
  as a response to the caller. (Notes: the calls to `https://gateway.marvel.com/v1/public/characters` are to account for
  the fact that there were new Marvel character(s) added where are not in `ids.txt`)
- If it is not the first call after the api starts, the api returns `ids` as a response to the caller.

### 2. `ids.txt` does not exist

The api makes calls to `https://gateway.marvel.com/v1/public/characters` with `offset` starting with `0` (and increment
of 100 in subsequent calls) and `limit` of 100 until there is no character returned. In each call, the character ids are
appended to an integer array, `ids`. The file `ids.txt` is created, and `ids` is written to `ids.txt`. The `ids` is
returned as a response to the caller.

##### Sample content of `ids.txt`

``` text
1011334
1017100
1009144
1010699
1009146
1016823
1009148
1009149
1010903
...
```

## Future Enhancements

- Add caching of `id`, `name` and `description`, read from cache and return it as response to calls to `/characters/{characterId}` to
  reduce latency.
- Add `TLS/https` to encrypt requests and responses.
- Refactor codes for better readability and maintainability.


