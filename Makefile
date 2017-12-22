

PREFIX := /
VAR := var/
RUN := run/
LOG := log/
ETC := etc/
USR := usr/
LOCAL := local/
VERSION := 0.18


CC := musl-gcc
COMPILER := "-compiler gccgo"

COMPILER_FLAGS := '-ldflags \'-linkmode external -extldflags "-static" "-fPIE" "-pie"\''

build:
	go get github.com/eyedeekay/gosam
	go build -o bin/si-i2p-plugin ./src
	@echo 'built'

build-static:
	go get github.com/eyedeekay/gosam
	go build -ldflags '-linkmode external -extldflags "-static"' \
		-o bin/si-i2p-plugin-static \
		./src

build-gccgo-static:
	go get github.com/eyedeekay/gosam
	go build "$(COMPILER)" \
		-gccgoflags '-extldflags "-fPIE" "-static" "-pie"' \
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


debug: build
	gdb ./bin/si-i2p-plugin

run:
	./bin/si-i2p-plugin 2>err | tee log

try: build
	./bin/si-i2p-plugin >log 2>err &
	sleep 1
	tail -f log

memcheck: build
	valgrind ./bin/si-i2p-plugin 2>err >log 2>err &
	tail -f log

test:	test-easy test-hard test-real test-diff test-dfhd test-dfsd # test-fake test-less test-loop test-fuzz

test-fake:
	@echo " It should not simply crash upon recieving a bad request, instead"
	@echo "it should log it, not make the request or touch the network at"
	@echo "all, and move on."
	echo http://notarealurl.i2p > parent/send
	cat parent/recv

test-less:
	@echo " It should not simply crash upon recieving a bad request, instead"
	@echo "it should log it, not make the request or touch the network at"
	@echo "all, and move on. This should include urls that don't exist under"
	@echo "domains that do."
	echo http://i2p-projekt.i2p/download > parent/send
	cat parent/recv

test-easy:
	@echo " It should know how to send requests for well-formed http url's"
	@echo "that point to b32 addresses or sites in the address book"
	echo http://i2p-projekt.i2p > parent/send
	#cat parent/recv
	cat i2p-projekt.i2p/recv

test-hard:
	@echo " It should also be able to recognize and correct simple"
	@echo "formatting mistakes in URL's and correct them where appropriate."
	echo i2p-projekt.i2p > parent/send
	#cat parent/recv
	cat i2p-projekt.i2p/recv

test-real:
	@echo " It should also be able to recognize and correct simple"
	@echo "formatting mistakes in URL's and correct them where appropriate."
	echo i2p-projekt.i2p/en/download > parent/send
	#cat parent/recv
	cat i2p-projekt.i2p/en/download/recv

test-df: test-diff test-dfhd test-dfsd

test-diff:
	@echo " It should know how to send requests for well-formed http url's"
	@echo "that point to b32 addresses or sites in the address book"
	echo http://inr.i2p > parent/send
	#cat parent/recv
	cat inr.i2p/recv

test-dfhd:
	@echo " It should know how to send requests for well-formed http url's"
	@echo "that point to b32 addresses or sites in the address book"
	echo inr.i2p > parent/send
	#cat parent/recv
	cat inr.i2p/recv

test-dfsd:
	@echo " It should know how to send requests for well-formed http url's"
	@echo "that point to b32 addresses or sites in the address book"
	echo inr.i2p/latest > parent/send
	#cat parent/recv
	cat inr.i2p/latest/recv

test-loop:
	@echo " It's rude and a privacy risk to use i2p-projekt.i2p(or any" "
	@echo "nonlocal eepSite as a test url in the final application. Instead,"
	@echo "it should generate a destination on this machine and query it as"
	@echo "a test, then immediately tear down the test tunnel."
	echo test.i2p > parent/serv
	echo http://test.i2p > parent/send


clean:
	killall si-i2p-plugin; \
	rm -rf parent *.i2p bin/si-i2p-plugin bin/si-i2p-plugin-static *.html *-pak *err *log static-include static-exclude del recv

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

kitten:
	cat i2p-projekt.i2p/recv

test-cat:
	cat i2p-projekt.i2p/recv > recv.html

exit:
	echo y > parent/del

noexit:
	echo n > parent/del

user:
	adduser --system --no-create-home --disabled-password --disabled-login --group sii2pplugin

checkinstall: build postinstall-pak postremove-pak description-pak
	checkinstall --default \
		--install=no \
		--fstrans=yes \
		--maintainer=problemsolver@openmailbox.org \
		--pkgname="si-i2p-plugin" \
		--pkgversion="$(VERSION)" \
		--pkglicense=gpl \
		--pkggroup=net \
		--pkgsource=./src/ \
		--deldoc=yes \
		--deldesc=yes \
		--delspec=yes \
		--backup=no \
		--pakdir=../

checkinstall-static: build postinstall-pak postremove-pak description-pak static-include static-exclude
	make static
	checkinstall --default \
		--install=no \
		--fstrans=yes \
		--maintainer=problemsolver@openmailbox.org \
		--pkgname="si-i2p-plugin" \
		--pkgversion="$(VERSION)-static" \
		--pkglicense=gpl \
		--pkggroup=net \
		--pkgsource=./src/ \
		--deldoc=yes \
		--deldesc=yes \
		--delspec=yes \
		--backup=no \
		--exclude=static-exclude \
		--include=static-include \
		--pakdir=../

postinstall-pak:
	@echo "#! /bin/sh" | tee postinstall-pak
	@echo "adduser --system --no-create-home --disabled-password --disabled-login --group sii2pplugin; true" | tee -a postinstall-pak
	@echo "mkdir -p $(PREFIX)$(VAR)$(LOG)si-i2p-plugin/ $(PREFIX)$(VAR)$(RUN)si-i2p-plugin/ || exit 1" | tee -a postinstall-pak
	@echo "chown -R sii2pplugin:adm $(PREFIX)$(VAR)$(LOG)si-i2p-plugin/ $(PREFIX)$(VAR)$(RUN)si-i2p-plugin/ || exit 1" | tee -a postinstall-pak
	@echo "exit 0" | tee -a postinstall-pak
	chmod +x postinstall-pak

postremove-pak:
	@echo "#! /bin/sh" | tee postremove-pak
	@echo "deluser sii2pplugin; true" | tee -a postremove-pak
	@echo "exit 0" | tee -a postremove-pak
	chmod +x postremove-pak

description-pak:
	@echo "si-i2p-plugin" | tee description-pak
	@echo "" | tee -a description-pak
	@echo "Destination-isolating http proxy for i2p. Keeps multiple eepSites" | tee -a description-pak
	@echo "from sharing a single reply destination, to limit the use of i2p" | tee -a description-pak
	@echo "metadata for fingerprinting purposes" | tee -a description-pak

static-include:
	@echo 'bin/si-i2p-plugin-static /usr/local/bin/' | tee static-include

static-exclude:
	@echo 'bin/si-i2p-plugin' | tee static-exclude


static:
	docker rm -f si-i2p-plugin-static; true
	docker build --force-rm -f Dockerfiles/Dockerfile.static -t si-i2p-plugin-static .
	docker run --name si-i2p-plugin-static -t si-i2p-plugin-static
	docker cp si-i2p-plugin-static:/opt/bin/si-i2p-plugin-static ./bin/si-i2p-plugin-static

uuser:
	docker build --force-rm -f Dockerfiles/Dockerfile.uuser -t si-i2p-plugin-uuser .
	docker run -d --rm --name si-i2p-plugin-uuser -t si-i2p-plugin-uuser
	docker exec -t si-i2p-plugin-uuser tail -n 1 /etc/passwd | tee si-i2p-plugin/passwd
	docker cp si-i2p-plugin-uuser:/bin/bash-static si-i2p-plugin/bash
	docker cp si-i2p-plugin-uuser:/bin/busybox si-i2p-plugin/busybox
	docker rm -f si-i2p-plugin-uuser; docker rmi -f si-i2p-plugin-uuser

docker:
	make static
	make uuser
	docker build --force-rm -f Dockerfiles/Dockerfile -t si-i2p-plugin .

docker-run:
	docker run -d \
		--cap-drop all \
		--name si-i2p-plugin \
		--user sii2pplugindocker \
		-t si-i2p-plugin

mps:
	bash -c "ps aux | grep si-i2p-plugin | grep -v gdb |  grep -v grep | grep -v https" 2>/dev/null

mls:
	ls -R *.i2p 2>/dev/null

ls:
	while true; do make -s mls 2>/dev/null; sleep 2; clear; done

ps:
	while true; do make -s mps 2>/dev/null; sleep 2; clear; done

