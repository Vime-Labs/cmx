VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  := -ldflags "-X github.com/Vime-Labs/cmx/cmd.Version=$(VERSION) -s -w"
DIST     := dist

.PHONY: build build-all test clean

## build: compila para a plataforma atual
build:
	go build $(LDFLAGS) -o cmx .

## build-all: compila para todas as plataformas em dist/
build-all: clean
	@mkdir -p $(DIST)
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o $(DIST)/cmx-linux-amd64     .
	GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o $(DIST)/cmx-linux-arm64     .
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o $(DIST)/cmx-darwin-amd64    .
	GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o $(DIST)/cmx-darwin-arm64    .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST)/cmx-windows-amd64.exe .
	@echo "\nBinários em $(DIST)/:"
	@ls -lh $(DIST)/

## test: roda todos os testes
test:
	go test ./...

## clean: remove binários gerados
clean:
	rm -rf $(DIST) cmx cmx.exe
