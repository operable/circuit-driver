BUILD_DIR    := _build
EXE_FILE     := $(BUILD_DIR)/circuit-driver
DOCKER_IMAGE ?= "operable/circuit-driver:dev"

# protobuf tooling
PROTOC_BIN         := $(shell which protoc)
PROTOC_DIR         := $(dir $(PROTOC_BIN))
PROTO_ROOT         := $(abspath $(addsuffix .., $(addprefix $(PROTOC_DIR), $(dir $(shell readlink -n $(PROTOC_BIN))))))
PROTO_ROOT_INCLUDE := $(addsuffix /include/, $(PROTO_ROOT))

# TODO This only works if GOPATH is a single directory
GOBIN_DIR          := $(addsuffix /bin, $(shell go env GOPATH))
GOFAST_PROTOC_BIN  := $(GOBIN_DIR)/protoc-gen-gofast

PROTO_DEFS  := $(wildcard api/*.proto)
PROTO_IMPLS := $(patsubst %.proto,%.pb.go,$(PROTO_DEFS))

.PHONY: all test exe clean docker vet deps pb-clean

all: test exe

deps:
	govendor sync

vet:
	govendor vet -x +local

$(PROTO_IMPLS): $(PROTO_DEFS)
	$(PROTOC_BIN) --plugin=$(GOFAST_PROTOC_BIN) --proto_path=$(PROTO_ROOT_INCLUDE):vendor:api --gofast_out=api $^

test: $(PROTO_IMPLS)
	govendor test +local -cover

# This is only intended to run in Travis CI and requires goveralls to
# be installed.
ci-coveralls: $(PROTO_IMPLS)
	goveralls -service=travis-ci

exe: $(PROTO_IMPLS) | $(BUILD_DIR)
	govendor build -o $(EXE_FILE)

docker:
	make clean
	GOOS=linux GOARCH=amd64 make exe
	docker build -t $(DOCKER_IMAGE) .

clean:
	rm -rf $(BUILD_DIR)

pb-clean:
	rm -f $(PROTO_IMPLS)

$(BUILD_DIR):
	mkdir -p $@
