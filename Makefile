GO_VERSION ?= 1.15
DOCKER_IMAGE ?= ghcr.io/tinyzimmer/go-gst:$(GO_VERSION)
GOLANGCI_VERSION ?= 1.31.0
GOLANGCI_LINT ?= _bin/golangci-lint
GOLANGCI_DOWNLOAD_URL ?= https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_VERSION}/golangci-lint-${GOLANGCI_VERSION}-$(shell uname | tr A-Z a-z)-amd64.tar.gz

$(GOLANGCI_LINT):
	mkdir -p $(dir $(GOLANGCI_LINT))
	cd $(dir $(GOLANGCI_LINT)) && curl -JL $(GOLANGCI_DOWNLOAD_URL) | tar xzf -
	chmod +x $(dir $(GOLANGCI_LINT))golangci-lint-$(GOLANGCI_VERSION)-$(shell uname | tr A-Z a-z)-amd64/golangci-lint
	ln -s golangci-lint-$(GOLANGCI_VERSION)-$(shell uname | tr A-Z a-z)-amd64/golangci-lint $(GOLANGCI_LINT)

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
docker-run: docker-build
	docker run --rm -it --privileged \
	    -v /lib/modules:/lib/modules:ro \
	    -v /sys:/sys:ro \
	    -v /usr/src:/usr/src:ro \
	    -v "$(PWD)":/workspace \
		-e HOME=/tmp \
	    $(DOCKER_IMAGE) $(CMD)

docker-lint: docker-build
	$(MAKE) docker-run CMD="make lint"