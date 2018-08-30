
UNAME ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
UARCH ?= $(shell uname -m | tr '[:upper:]' '[:lower:]' | sed 's|x86_64|amd64|g')

i2pd_dat?=i2pd_dat
browser=$(PWD)/browser

UPDATE_URL=https://www.torproject.org/projects/torbrowser/RecommendedTBBVersions

#COMMENT THE FOLLOWING LINE IF YOU WANT TO USE THE EXPERIMENTAL TBB
BROWSER_VERSION = $(shell curl $(UPDATE_URL) 2> /dev/null | grep -vi macos | grep -vi windows | grep -vi linux | head -n 2 | tail -n 1 | tr -d '",')

#UNCOMMENT THE FOLLOWING LINES IF YOU WANT TO USE THE EXPERIMENTAL TBB
BROWSER_VERSION=8.0a10

different_port=7073
DISPLAY = :0
BROWSER_PORT=44443
HOST=172.80.80.4

PREFIX := /
VAR := var/
RUN := run/
LIB := lib/
LOG := log/
ETC := etc/
USR := usr/
LOCAL := local/
VERSION := 0.20

OUTFOLDER = $(PWD)/bin

GOPATH = $(PWD)/.go

CGO_ENABLED=0

GO_COMPILER_OPTS = -a -tags netgo -ldflags '-w -extldflags "-static"'

info:
	@echo "Version $(VERSION)"
	@echo "$(UNAME), $(UARCH)"
	@echo "$(BROWSER_VERSION)"
	@echo "$(GOPATH)"

include misc/Makefiles/docker.mk
include misc/Makefiles/checkinstall.mk

rebuild: clean build

build: bin/si-i2p-plugin

nodeps: clean
	cd ./src/main/ && \
	GOOS=linux GOARCH=amd64 go build \
		$(GO_COMPILER_OPTS) \
		-o $(OUTFOLDER)/si-i2p-plugin \
		./si-i2p-plugin.go
	@echo 'built'

deps:
	go get -u github.com/eyedeekay/samrtc/src
	go get -u github.com/eyedeekay/jumphelper/src
	go get -u github.com/eyedeekay/gosam
	go get -u github.com/armon/go-socks5
	go get -u github.com/eyedeekay/si-i2p-plugin/src/errors
	go get -u github.com/eyedeekay/si-i2p-plugin/src/addresshelper
	go get -u github.com/eyedeekay/si-i2p-plugin/src/helpers
	go get -u github.com/eyedeekay/si-i2p-plugin/src/client
	go get -u github.com/eyedeekay/si-i2p-plugin/src/resolver
	go get -u github.com/eyedeekay/si-i2p-plugin/src/server
	go get -u github.com/eyedeekay/si-i2p-plugin/src
	#go get -u crawshaw.io/littleboss

bin/si-i2p-plugin:
	cd ./src/main/ && \
	GOOS=linux GOARCH=amd64 go build \
		$(GO_COMPILER_OPTS) \
		-o $(OUTFOLDER)/si-i2p-plugin \
		./si-i2p-plugin.go
	@echo 'built'

bin/si-i2p-plugin.app:
	cd ./src/main/ && \
	GOOS=darwin GOARCH=amd64 go build \
		$(GO_COMPILER_OPTS) \
		-o $(OUTFOLDER)/si-i2p-plugin.app \
		./si-i2p-plugin.go
	@echo 'built'

osx: bin/si-i2p-plugin.app

dmg:

bin/si-i2p-plugin.exe:
	cd ./src/main/ && \
	GOOS=windows GOARCH=amd64 go build \
		$(GO_COMPILER_OPTS) \
		-buildmode=exe \
		-o $(OUTFOLDER)/si-i2p-plugin.exe \
		./si-i2p-plugin.go
	@echo 'built'

windows: bin/si-i2p-plugin.exe

bin: bin/si-i2p-plugin bin/si-i2p-plugin.app bin/si-i2p-plugin.exe

build-arm: bin/si-i2p-plugin-arm

bin/si-i2p-plugin-arm: arm

noopts:
	cd ./src/main/ && \
	go build \
		-o $(OUTFOLDER)/si-i2p-plugin \
		./si-i2p-plugin.go
	@echo 'built'

arm:
	cd ./src/main/ && \
	ARCH=arm GOARCH=arm GOARM=7 go build \
		-compiler gc \
		$(GO_COMPILER_OPTS) \
		-buildmode=pie \
		-o $(OUTFOLDER)/si-i2p-plugin-arm \
		./si-i2p-plugin.go
	@echo 'built'

release: deps
	cd ./src/main/ && \
	GOOS="$(UNAME)" GOARCH="$(UARCH)" go build \
		$(GO_COMPILER_OPTS) \
		-buildmode=pie \
		-o $(OUTFOLDER)/si-i2p-plugin \
		./si-i2p-plugin.go
	@echo 'built release'

native:
	cd ./src/main/ && \
	go build \
		-a \
		-buildmode=pie \
		-o $(OUTFOLDER)/si-i2p-plugin \
		./si-i2p-plugin.go
	@echo 'built release'

android: bin/si-i2p-plugin-arm-droid

bin/si-i2p-plugin-arm-droid:
	cd ./src/main/ && \
	gomobile build \
		-target=android \
		$(GO_COMPILER_OPTS) \
		-o $(OUTFOLDER)/si-i2p-plugin-droid \
		./src/android/si-i2p-plugin.go
	@echo 'built'

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
	install -d -g sii2pplugin -o sii2pplugin -m744 $(PREFIX)$(VAR)$(LOG)/si-i2p-plugin/ $(PREFIX)$(VAR)$(LIB)/si-i2p-plugin/ $(PREFIX)$(ETC)si-i2p-plugin/
	install -d -g sii2pplugin -o sii2pplugin -m700 $(PREFIX)$(VAR)$(RUN)si-i2p-plugin/
	install -D -m755 bin/si-i2p-plugin $(PREFIX)$(USR)$(LOCAL)/bin/
	install -D -m755 $(ETC)init.d/si-i2p-plugin $(PREFIX)$(ETC)init.d/
	install -D -m644 $(ETC)systemd/sii2pplugin.service $(PREFIX)$(ETC)systemd/system/
	install -D -m644 $(ETC)apparmor.d/usr.bin.si-i2p-plugin $(PREFIX)$(ETC)apparmor.d/
	cp -r $(VAR)$(LIB)/si-i2p-plugin/ $(PREFIX)$(VAR)$(LIB)/si-i2p-plugin/
	install -D -g sii2pplugin -o sii2pplugin -m644 $(ETC)si-i2p-plugin/settings.cfg $(PREFIX)$(ETC)si-i2p-plugin/
	install -D -g sii2pplugin -o sii2pplugin -m600 $(ETC)si-i2p-plugin/addresses.csv $(PREFIX)$(ETC)si-i2p-plugin/

remove:
	rm -f $(PREFIX)$(USR)$(LOCAL)/bin/si-i2p-plugin \
		$(PREFIX)$(USR)$(LOCAL)/bin/si-i2p-plugin.sh \
		$(PREFIX)$(ETC)init.d/si-i2p-plugin \
		$(PREFIX)$(ETC)apparmor.d/usr.bin.si-i2p-plugin \
		$(PREFIX)$(ETC)apparmor.d/usr.local.bin.si-i2p-plugin \
		$(PREFIX)$(ETC)systemd/system/sii2pplugin.service \
		$(PREFIX)$(ETC)si-i2p-plugin/settings.cfg
	rm -rf $(PREFIX)$(VAR)$(LOG)/si-i2p-plugin/ $(PREFIX)$(VAR)$(RUN)si-i2p-plugin/ $(PREFIX)$(ETC)si-i2p-plugin/ $(PREFIX)$(VAR)$(LIB)/si-i2p-plugin/

run: nodeps
	./bin/si-i2p-plugin -proxy-port="$(BROWSER_PORT)" -addresshelper='http://inr.i2p,http://stats.i2p' 2>&1 | tee run.log

verbose: nodeps
	./bin/si-i2p-plugin -proxy-port="$(BROWSER_PORT)" -verbose=true -addresshelper='http://inr.i2p,http://stats.i2p' 2>&1 | tee run.log

try: nodeps
	./bin/si-i2p-plugin -proxy-port="$(BROWSER_PORT)" -conn-debug=true -addresshelper='http://inr.i2p,http://stats.i2p' 2>&1 | tee run.log

follow:
	docker logs -f si-proxy

clean:
	rm -rf parent services ./.*.i2p*/ ./*.i2p*/ \
		*.html *-pak *err *log \
		static-include static-exclude \
		bin/si-i2p-plugin* bin/si-i2p-plugin-arm var/lib/* \
		src/client/base64 src/client/id src/client/name \
		src/client/recv src/client/del src/client/send src/client/time \
		test/ src/test/ src/*/test/

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
	while true; do make docker-setup docker-run; sleep 30m; done

c: continuously

search:
	surf https://trac.torproject.org/projects/tor/ticket/25564

gotest:
	cd src && go test && \
		cd addresshelper && go test && \
		cd ../errors && go test && \
		cd ../helpers && go test && \
		cd ../client && go test && \
		cd ../server && go test && \
		cd ../main && go test

golint:
	cd src && golint && \
		cd addresshelper && golint && \
		cd ../errors && golint && \
		cd ../helpers && golint && \
		cd ../client && golint && \
		cd ../server && golint && \
		cd ../main && golint

govet:
	cd src && go vet && \
		cd addresshelper && go vet && \
		cd ../errors && go vet && \
		cd ../helpers && go vet && \
		cd ../client && go vet && \
		cd ../server && go vet && \
		cd ../main && go vet

gofmt:
	gofmt -w ./src/*.go ./src/*/*.go

gostuff: gofmt golint govet gotest

jhc:
	docker rm -f sam-jumphelper; docker rmi -f eyedeekay/sam-jumphelper

thing:
	http_proxy=http://127.0.0.1:44443 surf bn3x2dtvxov6jyfywc4oe2eoucktl7qhc45oomvr3vyvqtplce6a.b32.i2p.b32.i2p
