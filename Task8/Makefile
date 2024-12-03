all:
	make build

build:
	go build -o xkcd ./cmd/xkcd

server:
	go build -o xkcd-server ./cmd/xkcd-server

test:
	@go test ./... -v -race -cover -coverprofile coverage/coverage.out ## TODO: -race
	@go tool cover -html coverage/coverage.out -o coverage/coverage.html
	rm pkg/*/*.db

test_get:
	curl -X GET localhost:8080/pics?search=call?trying?but?it?says?number?blocked

test_post:
	curl -X POST localhost:8080/update

test_login:
	curl -X POST localhost:8080/login -d '{"email": "user1@example.com", "password": "password1"}'

test_login_admin:
	curl -X POST localhost:8080/login -d '{"email": "admin@example.com", "password": "password3"}'


test_rate:
	curl -X GET localhost:8080/pics?search=call?trying?but?it?says?number?blocked
	curl -X GET localhost:8080/pics?search=call?trying?but?it?says?number?blocked
	curl -X GET localhost:8080/pics?search=call?trying?but?it?says?number?blocked
	curl -X GET localhost:8080/pics?search=call?trying?but?it?says?number?blocked
	curl -X GET localhost:8080/pics?search=call?trying?but?it?says?number?blocked
	curl -X GET localhost:8080/pics?search=call?trying?but?it?says?number?blocked
	curl -X GET localhost:8080/pics?search=call?trying?but?it?says?number?blocked
	curl -X GET localhost:8080/pics?search=call?trying?but?it?says?number?blocked
	curl -X GET localhost:8080/pics?search=call?trying?but?it?says?number?blocked
	curl -X GET localhost:8080/pics?search=call?trying?but?it?says?number?blocked
	curl -X GET localhost:8080/pics?search=call?trying?but?it?says?number?blocked
	curl -X GET localhost:8080/pics?search=call?trying?but?it?says?number?blocked
	curl -X GET localhost:8080/pics?search=call?trying?but?it?says?number?blocked