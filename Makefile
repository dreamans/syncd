GO_CMD=go
SYNCD_BIN=syncd_bin
SYNCD_BIN_PATH=./output/bin
SYNCD_ETC_PATH=./output/etc
SYNCD_PUBLIC_PATH=./output/public
SYNCD_LOG_PATH=./output/log
SYNCD_RES_PATH=./output/resource

.PHONY: all
all: clean build install

.PHONY: linux
linux: clean build-linux install

.PHONY: build
build:
	@echo "build syncd start >>>"
	GOPROXY=https://goproxy.io $(GO_CMD) mod tidy
	$(GO_CMD) build -o $(SYNCD_BIN) ./syncd/main.go
	@echo ">>> build syncd complete"

.PHONY: install
install:
	@echo "install syncd start >>>"
	mkdir -p $(SYNCD_BIN_PATH)
	mv $(SYNCD_BIN) $(SYNCD_BIN_PATH)/syncd
	mkdir -p $(SYNCD_ETC_PATH)
	cp ./syncd.ini $(SYNCD_ETC_PATH)
	cp -r ./public $(SYNCD_PUBLIC_PATH)
	mkdir -p $(SYNCD_LOG_PATH)
	cp -r ./resource $(SYNCD_RES_PATH)
	@echo ">>> install syncd complete"

.PHONY: clean
clean:
	@echo "clean start >>>"
	rm -fr ./output
	rm -f $(SYNCD_BIN)
	@echo ">>> clean complete"

.PHONY: build-linux
build-linux:
	@echo "build-linux start >>>"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_CMD) build -o $(SYNCD_BIN) ./syncd/main.go
	@echo ">>> build-linux complete"