GO_VERSION ?= 1.15
DOCKER_IMAGE ?= ghcr.io/tinyzimmer/go-gst:$(GO_VERSION)

GOPATH ?= $(shell go env GOPATH)
GOBIN  ?= $(GOPATH)/bin
GOLANGCI_VERSION ?= v1.33.0
GOLANGCI_LINT    ?= $(GOBIN)/golangci-lint

PLUGIN_GEN ?= "$(shell go env GOPATH)/bin/gst-plugin-gen"

$(GOLANGCI_LINT):
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) $(GOLANGCI_VERSION)

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run -v

docker-build:
	docker build . \
	    -f .github/Dockerfile \
	    --build-arg GO_VERSION=$(GO_VERSION) \
	    -t $(DOCKER_IMAGE)

docker-push: docker-build
	docker push $(DOCKER_IMAGE)

CMD ?= /bin/bash
docker-run:
	docker run --rm --privileged \
	    -v /lib/modules:/lib/modules:ro \
	    -v /sys:/sys:ro \
	    -v /usr/src:/usr/src:ro \
	    -v "$(PWD)":/workspace \
		-e HOME=/tmp \
	    $(DOCKER_IMAGE) $(CMD)

docker-lint:
	$(MAKE) docker-run CMD="make lint"

$(PLUGIN_GEN):
	cd cmd/gst-plugin-gen && go build -o $(PLUGIN_GEN) .

install-plugin-gen: $(PLUGIN_GEN)