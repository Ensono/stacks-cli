.PHONY: build install clean test integration dep release
VERSION=`egrep -o '[0-9]+\.[0-9a-z.\-]+' version.go`
GIT_SHA=`git rev-parse --short HEAD || echo`

### @GO111MODULE=on; go build -ldflags "-X main.GitSHA=${GIT_SHA:="latest"}" -o bin/amidostacks-scaffolding-cli .

build:
	@echo "Building amidostacks-scaffolding-cli..."
	@mkdir -p bin
	@GO111MODULE=on; go build -o bin/amidostacks-scaffolding-cli .

install:
	@echo "Installing amidostacks-scaffolding-cli..."
	@install -c bin/amidostacks-scaffolding-cli /usr/local/bin/amidostacks-scaffolding-cli

clean:
	@rm -f bin/*

test:
	@echo "Running tests..."
	@go test `go list ./... | grep -v vendor/`

integration:
	@echo "Running integration tests..."
	@for i in `find ./integration -name test.sh`; do \
		echo "Running $$i"; \
		bash $$i || exit 1; \
		bash integration/expect/check.sh || exit 1; \
		rm /tmp/amidostacks-scaffolding-cli-*; \
	done

dep:
	@dep ensure

# TODO: update this as arm now on mac as well
release:
	@docker build -q -t amidostacks-scaffolding-cli_builder -f Dockerfile.build.alpine .
	@for platform in darwin linux windows; do \
		if [ $$platform == windows ]; then extension=.exe; fi; \
		docker run -it --rm -v ${PWD}:/app -e "GOOS=$$platform" -e "GOARCH=amd64" -e "CGO_ENABLED=0" amidostacks-scaffolding-cli_builder go build -ldflags="-s -w -X main.GitSHA=${GIT_SHA:=latest}" -o bin/amidostacks-scaffolding-cli-${VERSION}-$$platform-amd64$$extension; \
	done
	@docker run -it --rm -v ${PWD}:/app -e "GOOS=linux" -e "GOARCH=arm64" -e "CGO_ENABLED=0" amidostacks-scaffolding-cli_builder go build -ldflags="-s -w -X main.GitSHA=${GIT_SHA:=latest}" -o bin/amidostacks-scaffolding-cli-${VERSION}-linux-arm64;
