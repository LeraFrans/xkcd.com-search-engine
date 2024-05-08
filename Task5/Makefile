build:
	go build -o xkcd cmd/main.go

server:
	go build -o xkcd-server server/make_server.go

test_get:
	curl -X GET localhost:8080/pics?search=call?trying?but?it?says?number?blocked

test_post:
	curl -X POST -I localhost:8080/update

all: build