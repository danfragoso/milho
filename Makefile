git_revision != git log -1 --pretty="%h_%ad" --date=short
ldflags := -ldflags "-X 'github.com/danfragoso/milho.version=$(git_revision)'"

all: test

repl: 
	@go run $(ldflags) ./cli/*.go $(f)

build:
	@go build $(ldflags) -o milho ./cli/*.go

install:
	@cp milho /usr/bin/milho

wasm:
	@GOOS=js GOARCH=wasm go build $(ldflags) -o web/wasm/milho.wasm web/go/milho.go

test: test_tokenizer test_parser test_interpreter test_milho

test_tokenizer:
	@echo "Testing tokenizer...\n"
	@go test -v ./tokenizer
	@echo "\n--------"

test_parser:
	@echo "Testing parser...\n"
	@go test -v ./parser
	@echo "\n--------"

test_interpreter:
	@echo "Testing interpreter...\n"
	@go test -v ./interpreter
	@echo "\n--------"

test_milho: 
	@go test -v 
