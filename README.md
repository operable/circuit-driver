# circuit-driver

[![Build status](https://badge.buildkite.com/5e3d14a67525f02ccc31018c05fb47bb1952d1165904054dd5.svg?branch=master)](https://buildkite.com/operable/circuit-driver)

## Building Locally

### Prerequisites

```sh
brew install protobuf-c
go get -u github.com/gogo/protobuf/protoc-gen-gofast
go get -u github.com/kardianos/govendor
```

### Test and Build

```sh
make vet test exe
```

### Generate Docker Image

```
DOCKER_IMAGE=operable/circuit-driver:$VERSION make docker
```
