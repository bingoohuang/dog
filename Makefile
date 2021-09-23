.PHONY: default linux lint init
all: default

app=$(notdir $(shell pwd))
goVersion := $(shell go version | sed 's/go version //'|sed 's/ /_/')
buildTime := $(shell if hash gdate 2>/dev/null; then gdate --rfc-3339=seconds | sed 's/ /T/'; else date --rfc-3339=seconds | sed 's/ /T/'; fi)
# https://git-scm.com/docs/git-rev-list#Documentation/git-rev-list.txt-emaIem
gitCommit := $(shell git rev-list --oneline --format=format:'%h@%aI' --max-count=1 `git rev-parse HEAD` | tail -1)
#gitCommit := $(shell git rev-list -1 HEAD)
# https://stackoverflow.com/a/47510909
pkg := github.com/bingoohuang/gg/pkg/v
appVersion := 1.0.0
extldflags := -extldflags -static
# https://ms2008.github.io/2018/10/08/golang-build-version/
# https://github.com/kubermatic/kubeone/blob/master/Makefile
flags1 = "-s -w -X $(pkg).buildTime=$(buildTime) -X $(pkg).appVersion=$(appVersion) -X $(pkg).gitCommit=$(gitCommit) -X $(pkg).goVersion=$(goVersion)"
flags2 = "$(extldflags) -s -w -X $(pkg).buildTime=$(buildTime) -X $(pkg).appVersion=$(appVersion) -X $(pkg).gitCommit=$(gitCommit) -X $(pkg).goVersion=$(goVersion)"

init:

default:init
	go install -trimpath -ldflags=${flags1}  ./...
	upx ~/go/bin/${app}*

linux: init
	GOOS=linux GOARCH=amd64 go install -trimpath -ldflags=${flags2}  ./...
	upx ~/go/bin/linux_amd64/${app}*
	bssh scp ~/go/bin/linux_amd64/${app}* r:/usr/local/bin/
linux-arm64: init
	GOOS=linux GOARCH=arm64 go install -trimpath -ldflags=${flags2}  ./...
	upx ~/go/bin/linux_arm64/${app}*
	# ldd /usr/local/bin/busy
	bssh scp ~/go/bin/linux_amd64/${app}* r:/usr/local/bin/
lint-all:
	golangci-lint run --enable-all

lint:
	golangci-lint run ./...

fmt:
	gofumpt -l -w .
	gofmt -s -w .
	go mod tidy
	go fmt ./...
	revive .
	goimports -w .
	gci -w -local github.com/daixiang0/gci