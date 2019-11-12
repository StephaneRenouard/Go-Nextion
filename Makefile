COMPONENT=go-nextion

BINARIES=bin/$(COMPONENT)-armhf bin/screen-updater-armhf

.PHONY: $(BINARIES) clean

bin/$(COMPONENT)-armhf:
	env GOOS=linux GOARCH=arm go build -o $@

bin/screen-updater-armhf:
	env GOOS=linux GOARCH=arm go build -o $@ ./cmd/updater


all: $(BINARIES)

prepare:
	glide install

update:
	glide update

.ONESHELL:
deb-armhf: bin/$(COMPONENT)-armhf bin/screen-updater-armhf
	$(eval VERSION := $(shell cat ./version))
	$(eval ARCH := $(shell echo "armhf"))
	$(eval BUILD_NAME := $(shell echo "$(COMPONENT)-$(VERSION)-$(ARCH)"))
	$(eval BUILD_PATH := $(shell echo "build/$(BUILD_NAME)"))
	make deb VERSION=$(VERSION) BUILD_PATH=$(BUILD_PATH) ARCH=$(ARCH) BUILD_NAME=$(BUILD_NAME)

deb:
	mkdir -p $(BUILD_PATH)/usr/local/bin $(BUILD_PATH)/etc/systemd/system
	cp -r ./scripts/DEBIAN $(BUILD_PATH)/
	cp ./scripts/*.service $(BUILD_PATH)/etc/systemd/system/
	sed -i "s/amd64/$(ARCH)/g" $(BUILD_PATH)/DEBIAN/control
	sed -i "s/VERSION/$(VERSION)/g" $(BUILD_PATH)/DEBIAN/control
	sed -i "s/COMPONENT/$(COMPONENT)/g" $(BUILD_PATH)/DEBIAN/control
	cp ./scripts/Makefile $(BUILD_PATH)/../
	cp bin/screen-updater-$(ARCH) $(BUILD_PATH)/usr/local/bin/screen-updater
	mkdir -p $(BUILD_PATH)/data/screen
	cp cmd/update-nextion.sh $(BUILD_PATH)/usr/local/bin/update-nextion.sh
	chmod +x  $(BUILD_PATH)/usr/local/bin/update-nextion.sh
	cp bin/$(COMPONENT)-$(ARCH) $(BUILD_PATH)/usr/local/bin/$(COMPONENT)
	make -C build DEB_PACKAGE=$(BUILD_NAME) deb

clean:
	rm -rf bin build
