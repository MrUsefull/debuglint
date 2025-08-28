dirs:
	mkdir -p build/

test: dirs
	gotestsum --junitfile build/report.xml --format testname -- ./... -race -cover -covermode=atomic -coverprofile=build/cover.out --count=1
	go tool cover -func build/cover.out | grep total

htmlcover: test
	go tool cover -html build/cover.out -o build/cover.html

cideps:
	go install gotest.tools/gotestsum@latest
	go get github.com/boumenot/gocover-cobertura

ci: cideps test
	go run github.com/boumenot/gocover-cobertura < build/cover.out > build/coverage.xml

clean:
	rm -rf build

.PHONY: lint
lint:
	golangci-lint run 
