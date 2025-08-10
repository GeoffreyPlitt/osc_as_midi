DOCKER_IMAGE := osc-midi-bridge
CONTAINER := osc-midi-bridge-dev

.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

.PHONY: build
build: docker-build
	docker run --rm -v $(PWD):/app -v ~/go/pkg/mod:/go/pkg/mod $(DOCKER_IMAGE) go build -o osc-midi-bridge

.PHONY: test
test: docker-build
	docker run --rm -v $(PWD):/app -v ~/go/pkg/mod:/go/pkg/mod $(DOCKER_IMAGE) go test -v -cover ./...

.PHONY: integration-test
integration-test: build
	docker run --rm -v $(PWD):/app $(DOCKER_IMAGE) bash integration_test.sh

.PHONY: dev
dev: docker-build
	docker run --rm -it --name $(CONTAINER) \
		-v $(PWD):/app \
		-v ~/go/pkg/mod:/go/pkg/mod \
		-p 9000:9000 \
		-e DEBUG=$(DEBUG) \
		$(DOCKER_IMAGE) \
		bash -c "jackd -d dummy -r 48000 -p 64 & sleep 1 && bash"

.PHONY: run
run: build
	docker run --rm -it \
		-v $(PWD):/app \
		-p 9000:9000 \
		-e DEBUG=$(DEBUG) \
		$(DOCKER_IMAGE) \
		bash -c "jackd -d dummy -r 48000 -p 64 & sleep 1 && ./osc-midi-bridge $(ARGS)"

.PHONY: shell
shell:
	docker exec -it $(CONTAINER) /bin/bash

.PHONY: mod-download
mod-download: docker-build
	docker run --rm -v $(PWD):/app -v ~/go/pkg/mod:/go/pkg/mod $(DOCKER_IMAGE) go mod download

.PHONY: mod-tidy
mod-tidy: docker-build
	docker run --rm -v $(PWD):/app -v ~/go/pkg/mod:/go/pkg/mod $(DOCKER_IMAGE) go mod tidy

.PHONY: fmt
fmt: docker-build
	docker run --rm -v $(PWD):/app $(DOCKER_IMAGE) go fmt ./...

.PHONY: vet
vet: docker-build
	docker run --rm -v $(PWD):/app -v ~/go/pkg/mod:/go/pkg/mod $(DOCKER_IMAGE) go vet ./...

.PHONY: clean
clean:
	docker stop $(CONTAINER) 2>/dev/null || true
	docker rm $(CONTAINER) 2>/dev/null || true
	rm -f osc-midi-bridge