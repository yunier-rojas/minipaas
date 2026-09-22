.PHONY: test
test:
	@go test --cover -parallel=1 -v -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | sort -rnk3
