APP := vw
MODULE := github.com/comfucios/vw
CMD := ./cmd/vw
BIN_DIR := bin
DIST_DIR := dist

VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := \
	-s -w \
	-X $(MODULE)/internal/version.version=$(VERSION) \
	-X $(MODULE)/internal/version.commit=$(COMMIT) \
	-X $(MODULE)/internal/version.date=$(DATE)

.PHONY: all
all: build

.PHONY: build
build:
	mkdir -p $(BIN_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(APP) $(CMD)

.PHONY: run
run:
	go run $(CMD)

.PHONY: test
test:
	go test ./...

.PHONY: test-race
test-race:
	go test -race ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	rm -rf $(BIN_DIR) $(DIST_DIR)

.PHONY: doctor
doctor: build
	$(BIN_DIR)/$(APP) doctor

.PHONY: install
install:
	./scripts/install.sh

.PHONY: uninstall
uninstall:
	./scripts/uninstall.sh

.PHONY: package-with-bw
package-with-bw: build
	./scripts/package-with-bw.sh

.PHONY: snapshot
snapshot:
	goreleaser release --snapshot --clean

.PHONY: release-check
release-check:
	goreleaser check
