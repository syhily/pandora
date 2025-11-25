PREFIX ?= /opt/homebrew/bin

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

clean: ## Clean up build files.
	rm -rf pandora

deps: ## Update vendor.
	go mod verify
	go mod tidy -v
	go get -u ./...

fmt: ## Format code with goimports
	@command -v goimports >/dev/null 2>&1 || { echo "goimports not found, installing..."; go install golang.org/x/tools/cmd/goimports@latest; }
	goimports -local github.com/syhily/pandora -w .

build: clean ## Build executable files
	go build -o pandora main.go

gen: ## Generate Protobuf files
	protoc -I=internal/fonts/proto --go_out=internal/fonts/proto --go_opt=paths=source_relative internal/fonts/proto/index.proto

install: build ## Build and install this tool
	@echo "Installing pandora to $(PREFIX)"
	@mkdir -p $(PREFIX)
	mv pandora $(PREFIX)
	pandora completion zsh > ~/.zsh/completions/_pandora.zsh
