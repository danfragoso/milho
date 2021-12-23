.PHONY: repl build install www

git_revision != git log -1 --pretty="%h_%ad" --date=short
ldflags := -ldflags "-X 'github.com/danfragoso/milho.version=$(git_revision)'"
pwd != pwd

all: test

repl: 
	@go run $(ldflags) ./cli/*.go $(f)

build:
	@go build $(ldflags) -o milho ./cli/*.go

install:
	@cp milho /usr/bin/milho

www:
	@GOOS=js GOARCH=wasm go build $(ldflags) -o www/wasm/milho.wasm www/go/milho.go

test_spec:
	@make build
	@git clone https://github.com/milho-lang/tests
	@cd tests && ./run $(pwd)/milho
	@rm -rf tests

test: test_tokenizer test_parser test_interpreter test_milho test_spec

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
