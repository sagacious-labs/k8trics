PROTO_OUT_DIR := ./pkg/protos
API_PROTO_FILE := v1alpha1/api/api.proto
BASE_PROTO_FILE := v1alpha1/base/base.proto
GO_MODULE := github.com/sagacious-labs/k8trics

.PHONY: run
run: compile-proto
	go run ./cmd/.

.PHONY: compile
compile: compile-proto
	go build -o ./bin/k8trics ./cmd/.

.PHONY: mod-download
mod-download:
	go mod download

.PHONY: compile-proto
compile-proto:
	mkdir -p ./pkg/protos && \
	protoc --go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative \
	-I protos protos/$(API_PROTO_FILE) protos/$(BASE_PROTO_FILE)

.PHONY: docker
docker:
	@test -n "$(VERSION)" || (echo "VERSION is a required variable for target \"docker\"" ; exit 1)
	DOCKER_BUILDKIT=1 docker build . -t utkarsh23/k8trics:v$(VERSION)
