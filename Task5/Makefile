all:
	make build

build:
	go build -o xkcd ./cmd/xkcd

server:
	go build -o xkcd-server ./cmd/xkcd-server

test_get:
	curl -X GET localhost:8080/pics?search=call?trying?but?it?says?number?blocked

test_post:
	curl -X POST localhost:8080/update