
UNAME ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
UARCH ?= $(shell uname -m | tr '[:upper:]' '[:lower:]' | sed 's|x86_64|amd64|g')

i2pd_dat?=$(PWD)/i2pd_dat
browser=$(PWD)/browser
BROWSER_VERSION ?= $(shell curl https://www.torproject.org/projects/torbrowser.html.en 2>&1 | grep '<th>GNU/Linux<br>' | sed 's|<th>GNU/Linux<br><em>(||g' | sed 's|)</em></th>||g' | tr -d ' ')
different_port=7078
DISPLAY = :0
BROWSER_PORT=44443
HOST=172.80.80.4

PREFIX := /
VAR := var/
RUN := run/
LOG := log/
ETC := etc/
USR := usr/
LOCAL := local/
VERSION := 0.20

GOPATH = $(PWD)/.go

GO_COMPILER_OPTS = -a -tags netgo -ldflags '-w -extldflags "-static"'

info: gofmt
	@echo "Version $(VERSION)"
	@echo "$(UNAME), $(UARCH)"
	@echo "$(BROWSER_VERSION)"
	@echo "$(GOPATH)"

include misc/Makefiles/docker.mk
include misc/Makefiles/checkinstall.mk

rebuild: clean build

build: bin/si-i2p-plugin

nodeps: clean
	GOOS=linux GOARCH=amd64 go build \
		$(GO_COMPILER_OPTS) \
		-o bin/si-i2p-plugin \
		./src/main/si-i2p-plugin.go
	@echo 'built'

deps:
	go get -u github.com/eyedeekay/jumphelper/src
	go get -u github.com/eyedeekay/gosam
	go get -u github.com/armon/go-socks5
	go get -u github.com/eyedeekay/si-i2p-plugin/src/errors
	go get -u github.com/eyedeekay/si-i2p-plugin/src/addresshelper
	go get -u github.com/eyedeekay/si-i2p-plugin/src


bin/si-i2p-plugin: deps
	GOOS=linux GOARCH=amd64 go build \
		$(GO_COMPILER_OPTS) \
		-o bin/si-i2p-plugin \
		./src/main/si-i2p-plugin.go
	@echo 'built'

bin/si-i2p-plugin.bin: deps
	GOOS=darwin GOARCH=amd64 go build \
		$(GO_COMPILER_OPTS) \
		-o bin/si-i2p-plugin.bin \
		./src/main/si-i2p-plugin.go
	@echo 'built'

osx: bin/si-i2p-plugin.bin

bin/si-i2p-plugin.exe: deps
	GOOS=windows GOARCH=amd64 go build \
		$(GO_COMPILER_OPTS) \
		-buildmode=exe \
		-o bin/si-i2p-plugin.exe \
		./src/main/si-i2p-plugin.go
	@echo 'built'

windows: bin/si-i2p-plugin.exe

bin: bin/si-i2p-plugin bin/si-i2p-plugin.bin bin/si-i2p-plugin.exe

build-arm: bin/si-i2p-plugin-arm

bin/si-i2p-plugin-arm: deps arm

arm:
	ARCH=arm GOARCH=arm GOARM=7 go build \
		-compiler gc \
		$(GO_COMPILER_OPTS) \
		-buildmode=pie \
		-o bin/si-i2p-plugin-arm \
		./src/main/si-i2p-plugin.go
	@echo 'built'

release: deps
	GOOS="$(UNAME)" GOARCH="$(UARCH)" go build \
		$(GO_COMPILER_OPTS) \
		-buildmode=pie \
		-o bin/si-i2p-plugin \
		./src/main/si-i2p-plugin.go
	@echo 'built release'

native: deps
	go build \
		-a \
		-buildmode=pie \
		-o bin/si-i2p-plugin \
		./src/main/si-i2p-plugin.go
	@echo 'built release'

android: bin/si-i2p-plugin-arm-droid

bin/si-i2p-plugin-arm-droid: deps
	gomobile build \
		-target=android \
		$(GO_COMPILER_OPTS) \
		-o bin/si-i2p-plugin-droid \
		./src/android/si-i2p-plugin.go
	@echo 'built'

lib/si-i2p-plugin.a:
	GOOS="$(UNAME)" GOARCH="$(UARCH)" go build \
		$(GO_COMPILER_OPTS) \
		-buildmode=pie \
		-buildmode=archive \
		-o lib/si-i2p-plugin.a \
		./src/
	file lib/si-i2p-plugin.a

slib: lib/si-i2p-plugin.a

lib/si-i2p-plugin.so:
	GOOS="$(UNAME)" GOARCH="$(UARCH)" go build \
		-buildmode=pie \
		-buildmode=c-shared \
		-o lib/si-i2p-plugin.so \
		./src/main/si-i2p-plugin.go
	file lib/si-i2p-plugin.so

dylib: lib/si-i2p-plugin.so

lib: slib dylib

xpi2p:

debug: rebuild
	$(HOME)/.go/bin/dlv exec ./bin/si-i2p-plugin

dlv: rebuild
	$(HOME)/.go/bin/dlv debug ./src/main

all:
	make clobber; \
	make release; \
	make build-arm; \
	make checkinstall; \
	make checkinstall-arm; \
	make docker
	make tidy

install:
	mkdir -p $(PREFIX)$(VAR)$(RUN)si-i2p-plugin/
	install -d -g sii2pplugin -o sii2pplugin -m744 $(PREFIX)$(VAR)$(LOG)/si-i2p-plugin/ $(PREFIX)$(ETC)si-i2p-plugin/
	install -d -g sii2pplugin -o sii2pplugin -m700 $(PREFIX)$(VAR)$(RUN)si-i2p-plugin/
	install -D -m755 bin/si-i2p-plugin $(PREFIX)$(USR)$(LOCAL)/bin/
	install -D -m755 -g sii2pplugin -o sii2pplugin -m744 bin/si-i2p-plugin.sh $(PREFIX)$(USR)$(LOCAL)/bin/
	install -D -m755 -g sii2pplugin -o sii2pplugin -m744 bin/si-i2p-plugin-status.sh $(PREFIX)$(USR)$(LOCAL)/bin/
	install -D -m755 $(ETC)init.d/si-i2p-plugin $(PREFIX)$(ETC)init.d/
	install -D -m755 $(ETC)systemd/sii2pplugin.service $(PREFIX)$(ETC)systemd/system/
	install -D -g sii2pplugin -o sii2pplugin -m644 $(ETC)si-i2p-plugin/settings.cfg $(PREFIX)$(ETC)si-i2p-plugin/
	install -D -g sii2pplugin -o sii2pplugin -m600 $(ETC)si-i2p-plugin/addresses.csv $(PREFIX)$(ETC)si-i2p-plugin/

remove:
	rm -f $(PREFIX)$(USR)$(LOCAL)/bin/si-i2p-plugin \
		$(PREFIX)$(USR)$(LOCAL)/bin/si-i2p-plugin.sh \
		$(PREFIX)$(ETC)init.d/si-i2p-plugin \
		$(ETC)systemd/system/sii2pplugin.service \
		$(PREFIX)$(ETC)si-i2p-plugin/settings.cfg
	rm -rf $(PREFIX)$(VAR)$(LOG)/si-i2p-plugin/ $(PREFIX)$(VAR)$(RUN)si-i2p-plugin/ $(PREFIX)$(ETC)si-i2p-plugin/

run: nodeps
	./bin/si-i2p-plugin -proxy-port="4443" -addresshelper='http://inr.i2p,http://stats.i2p' 2>&1 | tee run.log

verbose: nodeps
	./bin/si-i2p-plugin -proxy-port="4443" -verbose=true -addresshelper='http://inr.i2p,http://stats.i2p' 2>&1 | tee run.log

try: nodeps
	./bin/si-i2p-plugin -proxy-port="4443" -conn-debug=true -addresshelper='http://inr.i2p,http://stats.i2p' 2>&1 | tee run.log

follow:
	docker logs -f si-proxy

clean:
	rm -rf parent services ./.*.i2p*/ ./*.i2p*/ \
		*.html *-pak *err *log \
		static-include static-exclude \
		bin/si-i2p-plugin bin/si-i2p-plugin-arm lib/*

kill:
	killall si-i2p-plugin; \
	rm -rf parent *.i2p parent

tidy:
	rm -rf parent *.i2p *.html *-pak *err *log static-include static-exclude

clobber: clean docker-clean
	rm -rf ../si-i2p-plugin_$(VERSION)*-1_amd64.deb

cat:
	cat parent/recv

exit:
	echo y > parent/del

noexit:
	echo n > parent/del

user:
	adduser --system --no-create-home --disabled-password --disabled-login --group sii2pplugin

golist:
	go list -f '{{.GoFiles}}' ./src

mps:
	bash -c "ps aux | grep si-i2p-plugin | grep -v gdb |  grep -v grep | grep -v https" 2> /dev/null

mls:
	@echo pipes
	@echo ==================
	ls *.i2p/* parent 2>/dev/null
	@echo

ls:
	while true; do make -s mls 2>/dev/null; sleep 2; clear; done

ps:
	while true; do make -s mps 2>/dev/null; sleep 2; clear; done

continuously:
	while true; do make docker-setup docker-run; sleep 10m; done

c: continuously

search:
	surf https://trac.torproject.org/projects/tor/ticket/25564

gotest:
	cd src && go test

golint:
	cd src && golint

govet:
	cd src && go vet

gofmt:
	gofmt -w ./src/*.go ./src/*/*.go
