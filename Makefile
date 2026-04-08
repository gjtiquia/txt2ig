.PHONY: all build web clean templ test

all: build

templ:
	templ generate

build: templ
	go build -o txt2ig .

web: templ
	go run . web --port 3000

test:
	go test ./...

clean:
	rm -f txt2ig
	find internal/web/templates -name "*.go" -delete