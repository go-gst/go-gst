
build-cmd:
	cd cmd/go-gst && go build -o ../../dist/go-gst

ARGS ?=
run-cmd: build-cmd
	dist/go-gst $(ARGS)

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
