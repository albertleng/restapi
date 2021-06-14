
# Marvel Characters Golang REST Api

A RESTful api that's based on [Marvel API](https://developer.marvel.com/) that returns:
- all the Marvel characters ids in a JSON array of numbers.
- the id, name and description of a character.

## Requirements
- Sign up for a free **Marvel Developer API** at https://developer.marvel.com/.
- [Go version 1.15](https://golang.org/dl/) or greater.
- Install [mux router](https://github.com/gorilla/mux) to implement a request router and dispatcher for matching incoming requests to their respective handler.
- Install [Go CORS handler](https://github.com/rs/cors) to handle cors requests.

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


