# Install air module for dev mode
.PHONY: air-install
air-install:
	go install github.com/air-verse/air@latest

# Install templ module
.PHONY: templ-install
templ-install:
	go install github.com/a-h/templ/cmd/templ@latest

# Generate templ
.PHONY: templ-gen
templ-gen:
	templ fmt ./templates
	templ generate

# Run the application
.PHONY: run
run:
	air

# Run all tests
.PHONY: test
test:
	go test -count=1 -p 1 ./...

# Check for direct dependency updates
.PHONY: check-updates
check-updates:
	go list -u -m -f '{{if not .Indirect}}{{.}}{{end}}' all | grep "\["
