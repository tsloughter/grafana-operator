OPERATOR_NAME  := grafana-operator
VERSION := $(shell date +%Y%m%d%H%M)
ACCOUNT := tsloughter
IMAGE := $(ACCOUNT)/$(OPERATOR_NAME)

.PHONY: install_deps build build-image clean

install_deps:
	glide install

build:
	rm -rf bin/%/$(OPERATOR_NAME)
	CGO_ENABLED=0 go build -v -i -o bin/$(OPERATOR_NAME) ./cmd

clean:
	rm -rf bin/%/$(OPERATOR_NAME)

bin/%/$(OPERATOR_NAME): clean
	GOOS=$* GOARCH=amd64 go build -v -i -o bin/$*/$(OPERATOR_NAME) ./cmd

build-image: bin/linux/$(OPERATOR_NAME)
	docker build . -t $(IMAGE):$(VERSION)

push-image: build-image
	docker push $(IMAGE):$(VERSION)