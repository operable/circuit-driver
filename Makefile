TOP_PKG                      = github.com/operable/circuit
PKG_DIRS                    := $(shell find . -not -path '*/\.*' -type d | grep -v _build | sort)
PKGS                        := $(TOP_PKG) $(subst ., $(TOP_PKG), $(PKG_DIRS))
BUILD_DIR                    = _build
EXE_FILE                    := $(BUILD_DIR)/circuit-driver

.PHONY: all test exe clean docker

all: Makefile test exe

test:
	go test -cover $(PKGS)

exe: $(EXE_FILE)

$(EXE_FILE): $(BUILD_DIR)
	go build -o $(BUILD_DIR)/circuit-driver github.com/operable/circuit-driver

clean:
	rm -rf $(BUILD_DIR)
	find . -name "*.test" -type f | xargs rm -f

$(BUILD_DIR):
	mkdir -p $@

docker:
	make clean
	GOOS=linux GOARCH=amd64 make exe
	docker build -t operable/circuit-driver .
