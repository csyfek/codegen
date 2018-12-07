.PHONY: clean charts all install qtsetup env-debian

BINARIES=structs2bindings structs2schema structs2interface codegengui db2structs

CGO_CXXFLAGS_ALLOW=".*"
CGO_LDFLAGS_ALLOW=".*"
CGO_CFLAGS_ALLOW=".*"
QT_DIR=/home/jackman/Qt5.11.1
QT_VERSION=5.11.1

all: $(BINARIES)

clean:
	rm -rf $(BINARIES)
	rm -rf ./cmd/codegengui/rcc*
	rm -rf ./cmd/codegengui/moc*
	rm -rf ./cmd/codegengui/deploy

codegengui:
	qtdeploy build desktop ./cmd/codegengui

structs2bindings:
	go build github.com/jackmanlabs/codegen/cmd/structs2bindings

structs2schema:
	go build github.com/jackmanlabs/codegen/cmd/structs2schema

structs2interface:
	go build github.com/jackmanlabs/codegen/cmd/structs2interface

install:
	go install github.com/jackmanlabs/codegen/cmd/structs2bindings
	go install github.com/jackmanlabs/codegen/cmd/structs2schema
	go install github.com/jackmanlabs/codegen/cmd/structs2interface
	go install github.com/jackmanlabs/codegen/cmd/codegengui

fmt:
	for d in $(shell go list -f {{.Dir}} ./...); do goimports -w $$d/*.go; done

QT_INSTALLER=qt-opensource-linux-x64-5.11.1.run
QT_DOWNLOAD_PATH=http://download.qt.io/archive/qt/5.11/5.11.1/$(QT_INSTALLER)
qtsetup:
	wget --continue $(QT_DOWNLOAD_PATH)
	chmod +x $(QT_INSTALLER)
	-./$(QT_INSTALLER)
	go get -v github.com/therecipe/qt/cmd/...
	$(GOPATH)/bin/qtsetup

env-debian:
	sudo apt-get -y install build-essential libglu1-mesa-dev libpulse-dev libglib2.0-dev