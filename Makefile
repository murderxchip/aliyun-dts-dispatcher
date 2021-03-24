.PHONY: test build deploy-master deploy-test clean build-staging build-master
VERSION=`egrep -o '[0-9]+\.[0-9a-z.\-]+' version.go`
DATETIME=`date +'%Y-%m-%d %H:%M:%S'`
GIT_SHA=`git rev-parse --short HEAD || echo`
# go cmd
GOBUILD = gox
OSARCH = linux/amd64
#path
SOURCE_FILE = main.go
#binary name
DTSENV = dev
DTS = dts-dispatcher
PROJECT = $(DTS)-$(DTSENV)
#remote path
REMOTE_PATH = /home/work/project/dts-dispatcher/

export GO111MODULE=on
export GOPROXY=https://mirrors.aliyun.com/goproxy/

test:
	go run main.go
local:
	@echo ${GIT_SHA}
	$(GOBUILD) -osarch="darwin/amd64" -ldflags "-X 'github.com/murderxchip/aliyun-dts-dispatcher/define.GitSHA=${GIT_SHA}'  -X 'github.com/murderxchip/aliyun-dts-dispatcher/define.GitDateTime=${DATETIME}'" -output=dts-dispatcher
build:
	mkdir -p bin/config
	@cp config/config.toml.${DTSENV} bin/config/config.toml
	$(GOBUILD) -osarch="$(OSARCH)" -ldflags "-X 'github.com/murderxchip/aliyun-dts-dispatcher/define.GitSHA=${GIT_SHA}'  -X 'github.com/murderxchip/aliyun-dts-dispatcher/define.GitDateTime=${DATETIME}'" -output=bin/$(PROJECT)
clean:
	rm -rf bin
	rm -rf dts-dispatcher
