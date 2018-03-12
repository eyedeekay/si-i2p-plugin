
PREFIX := /
VAR := var/
RUN := run/
LOG := log/
ETC := etc/
USR := usr/
LOCAL := local/
VERSION := 0.19


COMPILER := "-compiler gccgo"
GO_COMPILER := "-compiler gc"

build: clean bin/si-i2p-plugin

bin/si-i2p-plugin:
	go get -u github.com/eyedeekay/gosam
	go build "$(COMPILER)" \
		-o bin/si-i2p-plugin \
		./src
	@echo 'built'


release:
	go get -u github.com/eyedeekay/gosam
	go build "$(GO_COMPILER)" -buildmode=pie \
		-o bin/si-i2p-plugin \
		./src
	@echo 'built release'


debug: build
	gdb ./bin/si-i2p-plugin

build-static:
	go get github.com/eyedeekay/gosam
	go build "$(GO_COMPILER)" -buildmode=pie \
		-a -ldflags '-extldflags "-static"' \
		-o bin/si-i2p-plugin-static \
		./src

build-gccgo-static:
	go get github.com/eyedeekay/gosam
	go build "$(COMPILER)" \
		-gccgoflags -extldflags "-static" -buildmode=pie\
		-o bin/si-i2p-plugin-static \
		./src

all:
	make clobber; \
	make; \
	make static; \
	make checkinstall; \
	make checkinstall-static; \
	make docker
	make tidy

install:
	mkdir -p $(PREFIX)$(VAR)$(LOG)/si-i2p-plugin/ $(PREFIX)$(VAR)$(RUN)si-i2p-plugin/ $(PREFIX)$(ETC)si-i2p-plugin/
	install -D bin/si-i2p-plugin $(PREFIX)$(USR)$(LOCAL)/bin/
	install -D bin/si-i2p-plugin.sh $(PREFIX)$(USR)$(LOCAL)/bin/
	install -D init.d/si-i2p-plugin $(PREFIX)$(ETC)init.d/
	install -D systemd/sii2pplugin.service $(PREFIX)$(ETC)systemd/system/
	install -D si-i2p-plugin/settings.cfg $(PREFIX)$(ETC)si-i2p-plugin/

remove:
	rm -f $(PREFIX)$(USR)$(LOCAL)/bin/si-i2p-plugin \
		$(PREFIX)$(USR)$(LOCAL)/bin/si-i2p-plugin.sh \
		$(PREFIX)$(ETC)init.d/si-i2p-plugin $(PREFIX)\
		$(ETC)systemd/system/sii2pplugin.service \
		$(PREFIX)$(ETC)si-i2p-plugin/settings.cfg
	rm -rf $(PREFIX)$(VAR)$(LOG)/si-i2p-plugin/ $(PREFIX)$(VAR)$(RUN)si-i2p-plugin/ $(PREFIX)$(ETC)si-i2p-plugin/

run: build
	./bin/si-i2p-plugin >run.log 2>run.err &

follow:
	tail -f run.log run.err | nl

try: build
	./bin/si-i2p-plugin -conn-debug=true >log 2>err &
	sleep 1
	tail -f log | nl

memcheck: release
	valgrind --track-origins=yes ./bin/si-i2p-plugin -conn-debug=true 1>log 2>err &
	sleep 2
	tail -f log | nl

clean:
	killall si-i2p-plugin; \
	rm -rf parent ./.*.i2p/ *.i2p/ bin/si-i2p-plugin bin/si-i2p-plugin-static *.html *-pak *err *log static-include static-exclude del recv

kill:
	killall si-i2p-plugin; \
	rm -rf parent *.i2p parent

tidy:
	rm -rf parent *.i2p *.html *-pak *err *log static-include static-exclude

clobber:
	rm -rf ../si-i2p-plugin_$(VERSION)*-1_amd64.deb
	docker rmi -f si-i2p-plugin-static si-i2p-plugin; true
	docker rm -f si-i2p-plugin-static si-i2p-plugin; true
	make clean

cat:
	cat parent/recv

exit:
	echo y > parent/del

noexit:
	echo n > parent/del

user:
	adduser --system --no-create-home --disabled-password --disabled-login --group sii2pplugin

static:
	docker rm -f si-i2p-plugin-static; true
	docker build --force-rm -f Dockerfiles/Dockerfile.static -t eyedeekay/si-i2p-plugin-static .
	docker run --name si-i2p-plugin-static -t eyedeekay/si-i2p-plugin-static
	docker cp si-i2p-plugin-static:/opt/bin/si-i2p-plugin-static ./bin/si-i2p-plugin-static

uuser:
	docker build --force-rm -f Dockerfiles/Dockerfile.uuser -t eyedeekay/si-i2p-plugin-uuser .
	docker run -d --rm --name si-i2p-plugin-uuser -t eyedeekay/si-i2p-plugin-uuser
	docker exec -t si-i2p-plugin-uuser tail -n 1 /etc/passwd | tee si-i2p-plugin/passwd
	docker cp si-i2p-plugin-uuser:/bin/bash-static si-i2p-plugin/bash
	docker cp si-i2p-plugin-uuser:/bin/busybox si-i2p-plugin/busybox
	docker rm -f si-i2p-plugin-uuser; docker rmi -f eyedeekay/si-i2p-plugin-uuser

docker:
	make static
	make uuser
	docker build --force-rm -f Dockerfiles/Dockerfile -t eyedeekay/si-i2p-plugin .

docker-run:
	docker run \
		--cap-drop all \
		--name si-i2p-plugin \
		--user sii2pplugindocker \
		-p 44443:4443 \
		-t eyedeekay/si-i2p-plugin

docker-run-thirdeye:
	docker run \
		--name thirdeye-proxy \
		--network thirdeye \
		--network-alias thirdeye-proxy \
		--hostname thirdeye-proxy \
		--cap-drop all \
		--user sii2pplugindocker \
		-p 44443:4443 \
		-t eyedeekay/si-i2p-plugin

mps:
	bash -c "ps aux | grep si-i2p-plugin | grep -v gdb |  grep -v grep | grep -v https" 2>/dev/null

mls:
	@echo pipes
	@echo ==================
	ls *.i2p/* parent 2>/dev/null
	@echo

ls:
	while true; do make -s mls 2>/dev/null; sleep 2; clear; done

ps:
	while true; do make -s mps 2>/dev/null; sleep 2; clear; done

include misc/Makefiles/demo.mk
include misc/Makefiles/test.mk
include misc/Makefiles/checkinstall.mk
