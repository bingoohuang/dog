default:
	go install -trimpath -ldflags='-s -w' ./...
	upx ~/go/bin/dog
linux:
	GOOS=linux GOARCH=amd64 go install -trimpath -ldflags='-extldflags=-static -s -w' ./...
	upx ~/go/bin/linux_amd64/dog
	# ldd /usr/local/bin/busy
lint:
	go mod tidy
	gofmt -s -w .
	gofumpt -l -w .
	go fmt ./...
	revive .
	goimports -w .
	golangci-lint run --enable-all

