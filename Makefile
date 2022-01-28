
# install gofumpt: go install mvdan.cc/gofumpt@latest
GIT_VERSION?=$(shell git describe --always --tags)

gofmt:
	@GO111MODULE=off gofumpt -w -l $(shell find . -type f -name '*.go'| grep -v "/vendor/\|/.git/\|/git/\|.*_y.go")

test:
	go test

define build
		@GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
endef

pub_image:
	@echo $(GIT_VERSION)
	@docker buildx build --platform linux/amd64 \
    		-t pubrepo.jiagouyun.com/googleimages/shorturl:$(GIT_VERSION) . --push

local: gofmt test
	$(call build)

push: pub_image



