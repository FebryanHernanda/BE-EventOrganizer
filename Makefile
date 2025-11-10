.PHONY: docs
docs:
	swag init -g cmd/main.go -o docs

.PHONY: run
run:
	air
