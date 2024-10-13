BUF_VERSION:=1.5.0

.PHONY: build_frontend
build_frontend:
	cd frontend && npm run build

.PHONY: build
build:
	@echo "compiling source code"
	@go build -o ./bin/ -v ./cmd/...

.PHONY: test
test:
	@echo "running tests"
	@go test -v ./...

.PHONY: generate
generate:
	docker run -v $$(pwd):/src -w /src --rm bufbuild/buf:$(BUF_VERSION) generate

.PHONY: lint
lint:
	docker run -v $$(pwd):/src -w /src --rm bufbuild/buf:$(BUF_VERSION) lint
	docker run -v $$(pwd):/src -w /src --rm bufbuild/buf:$(BUF_VERSION) breaking --against 'https://github.com/johanbrandhorst/grpc-gateway-boilerplate.git#branch=master'


DOCKER_VERSION:=v1.0.2

.PHONY: docker_push
docker_push: build build_frontend
	docker buildx build --platform linux/amd64 -t inconetwork/cosmfaucet:$(DOCKER_VERSION) .
	docker push inconetwork/cosmfaucet:$(DOCKER_VERSION)
