build_all: build_lambda
	go build -o ./build/smh-webserver ./cmd/webserver/serve.go
	go build -o ./build/smh-configurator ./cmd/configurator/configure.go
	go build -o ./build/smh-runner ./cmd/runner/run.go
	go build -o ./build/rmq-direct-publisher ./cmd/rmq-proxy/publisher-direct/direct_publisher.go
	go build -o ./build/rmq-consumer ./cmd/rmq-proxy/consumer/consumer.go

build_lambda:
	GOOS=linux GOARCH=amd64 go build -o ./build/rmq-lambda-publisher ./cmd/rmq-proxy/publisher-lambda/lambda_publisher.go

test:
	go test ./...

coverage:
	go test -v -coverprofile ./build/cover.out ./...
	go tool cover -html=./build/cover.out -o ./build/cover.html