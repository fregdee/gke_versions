.PHONY: install
install:
	go install ./cmd/gke_versions

.PHONY: build
build:
	go build ./cmd/gke_versions

