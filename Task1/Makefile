




build:
	go build -o myapp main.go normalizator.go

run:
	go run main.go normalizator.go


compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o myapp-linux-arm main.go normalizator.go
	GOOS=linux GOARCH=arm64 go build -o myapp-linux-arm64 main.go normalizator.go
	GOOS=freebsd GOARCH=386 go build -o myapp-freebsd-386 main.go normalizator.go
	GOOS=windows GOARCH=amd64 go build -o myapp-windows-amd64.exe main.go normalizator.go

all: build