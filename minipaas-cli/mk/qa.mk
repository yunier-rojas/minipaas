.PHONY: qa
qa:
	@go fmt github.com/yunier-rojas/minipaas/minipaas-cli/internal/...
	@go fmt github.com/yunier-rojas/minipaas/minipaas-cli/cmd/...
	golangci-lint run


.PHONY: imports
imports:
	go run github.com/quantumcycle/go-import-checks@latest --config qa/.import.yaml
