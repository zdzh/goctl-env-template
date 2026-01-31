.PHONY: build install clean test

build-cli:
	go build -o goctl-env-template cmd/root.go

build-plugin:
	go build -o plugin/main plugin/main.go

build: build-cli build-plugin

install:
	go install github.com/zdzh/goctl-env-template/cmd@latest

clean:
	rm -f goctl-env-template
	rm -f plugin/main
	rm -f .env.template

test: build-cli
	./goctl-env-template -c config/config.go -o .env.template

demo: test
	@echo "=== Generated .env.template ==="
	@cat .env.template
