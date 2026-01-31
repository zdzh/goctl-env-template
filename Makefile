.PHONY: build install clean test demo

build:
	go build -o goctl-env-template .

install:
	go install github.com/zdzh/goctl-env-template@latest

clean:
	rm -f goctl-env-template
	rm -f .env.template

test:
	go test -v

demo: build
	./goctl-env-template -c config/config.go -o .env.template
	@echo "=== Generated .env.template ==="
	@cat .env.template
