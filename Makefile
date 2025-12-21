.PHONY: repl build install install-std www hooks

git_revision != git log -1 --pretty="%h_%ad" --date=short
ldflags := -ldflags "-X 'github.com/danfragoso/milho.version=$(git_revision)'"
pwd != pwd
home_dir := $(HOME)
std_path := $(home_dir)/.milho/core

all: test

repl:
	@go run $(ldflags) ./cli/*.go $(f)

build:
	@go build $(ldflags) -o milho ./cli/*.go

install:
	@cp milho /usr/local/bin/milho
	@echo "Milho installed to /usr/local/bin/milho"

install-std:
	@echo "Installing Milho standard library to $(std_path)..."
	@mkdir -p $(std_path)
	@cp core/*.milho $(std_path)/
	@echo ""
	@echo "✓ Standard library installed to $(std_path)"
	@echo ""
	@echo "Add this to your ~/.zshrc or ~/.bashrc:"
	@echo "  export MILHO_STD_PATH=$(std_path)"
	@echo ""
	@echo "Then imports like (import \"core/http.milho\") will work from anywhere!"

www:
	@GOOS=js GOARCH=wasm go build $(ldflags) -o www/wasm/milho.wasm www/go/milho.go

hooks:
	@cp ./hooks/pre-push.milho .git/hooks/pre-push
	@chmod +x .git/hooks/pre-push

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
