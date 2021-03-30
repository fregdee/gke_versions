VERSION=0.0.1
LDFLAGS=-ldflags "-w -s -X github.com/fregdee/gke_versions.version=${VERSION}"
GO111MODULE=on

.PHONY: all
all: gke_versions

.PHONY: gke_versions
gke_versions:
	go build $(LDFLAGS) -o gke_versions ./cmd/gke_versions

.PHONY: clean
clean:
	rm -rf gke_versions

.PHONY: check
check:
	go test ./...

.PHONY: tag
tag:
	git tag v${VERSION}
	git push origin v${VERSION}
	git push origin HEAD
