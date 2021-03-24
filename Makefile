all: test

run: 
	@go run ./repl/*.go

build:
	@go build -o milho ./repl/*.go

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