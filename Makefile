build-all: build-linux-amd64 build-win-amd64 build-darwin-amd64

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o build/xproxy-linux-amd64

build-win-amd64:
	GOOS=windows GOARCH=amd64 go build -o build/xproxy-win-amd64

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o build/xproxy-darwin-amd64