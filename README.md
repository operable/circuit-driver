# circuit-driver

[![Build Status](https://travis-ci.org/operable/circuit-driver.svg?branch=master)](https://travis-ci.org/operable/circuit-driver)

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
