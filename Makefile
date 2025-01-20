.PHONY: help build

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "    build   - build the project"
	@echo "    run     - run the project"
	@echo "    help    - display this help message"

build:
	@echo "Building the project..."
	@go build -o tmp/synapse
	@echo "Project built successfully"

run:
	./tmp/synapse
