

PREFIX := /
VAR := var/
RUN := run/
LOG := log/
ETC := etc/
USR := usr/
LOCAL := local/
VERSION := 0.1

all:
	go build -o bin/si-i2p-plugin ./src

install:
	mkdir -p $(PREFIX)$(VAR)$(LOG)/si-i2p-plugin/ $(PREFIX)$(VAR)$(RUN)si-i2p-plugin/ $(PREFIX)$(ETC)si-i2p-plugin/
	install bin/si-i2p-plugin $(PREFIX)$(USR)$(LOCAL)/bin/
	install bin/si-i2p-plugin.sh $(PREFIX)$(USR)$(LOCAL)/bin/
	install init.d/si-i2p-plugin $(PREFIX)$(ETC)init.d/
	install systemd/sii2pplugin.service $(PREFIX)$(ETC)systemd/system
	install si-i2p-plugin/settings.cfg $(PREFIX)$(ETC)si-i2p-plugin/settings.cfg

debug:
	make
	gdb ./bin/si-i2p-plugin

try:
	bash -c "./bin/si-i2p-plugin 1>log 2>err" & sleep 1 && true
	cat parent/recv | tee test.html

test:
	echo http://i2p-projekt.i2p > parent/send

alttest:
	./bin/si-i2p-plugin --url fireaxe.i2p

testother:
	echo http://i2p-projekt.i2p/en/download > parent/send

clean:
	rm -rf i2p-projekt.i2p bin/si-i2p-plugin *.html *-pak *err *log parent ../si-i2p-plugin_$(VERSION)-1_amd64.deb
	docker rmi -f si-i2p-plugin-static si-i2p-plugin-rpm si-i2p-plugin

cat:
	cat i2p-projekt.i2p/recv

name:
	cat i2p-projekt.i2p/name

exit:
	echo y > parent/del

noexit:
	echo n > parent/del

html:
	cat i2p-projekt.i2p/recv | tee test.html; true

htmlother:
	cat i2p-projekt.i2p/recv | tee test2.html; true

diff:
	diff test.html test2.html

html-test:
	sr W ./test.html

user:
	adduser --system --no-create-home --disabled-password --disabled-login --group sii2pplugin

checkinstall: all preinstall-pak postremove-pak description-pak
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

preinstall-pak:
	@echo "#! /bin/sh" | tee preinstall-pak
	@echo "adduser --system --no-create-home --disabled-password --disabled-login --group sii2pplugin || exit 1" | tee -a preinstall-pak
	@echo "chown -R sii2pplugin:adm $(PREFIX)$(VAR)$(LOG)/si-i2p-plugin/ $(PREFIX)$(VAR)$(RUN)si-i2p-plugin/ || exit 1" | tee -a preinstall-pak
	@echo "exit 0" | tee -a preinstall-pak
	chmod +x preinstall-pak

postremove-pak:
	@echo "#! /bin/sh" | tee postremove-pak
	@echo "deluser sii2pplugin || exit 1" | tee -a postremove-pak
	@echo "exit 0" | tee -a postremove-pak
	chmod +x postremove-pak

description-pak:
	@echo "si-i2p-plugin" | tee description-pak
	@echo "" | tee -a description-pak
	@echo "Destination-isolating http proxy for i2p. Keeps multiple eepSites" | tee -a description-pak
	@echo "from sharing a single reply destination, to limit the use of i2p" | tee -a description-pak
	@echo "metadata for fingerprinting purposes" | tee -a description-pak

static:
	docker build --force-rm -f Dockerfile/Dockerfile.static -t si-i2p-plugin-static .
	docker run -d --name si-i2p-plugin-static -t si-i2p-plugin-static
	docker cp si-i2p-plugin-static:/opt/bin/si-i2p-plugin ./bin/si-i2p-plugin-static
	docker rm -f si-i2p-plugin-static

docker:
	make static
	docker build --force-rm -f Dockerfile/Dockerfile -t si-i2p-plugin .
	docker rmi -f si-i2p-plugin-static

docker-run:
	docker run -d \
		--name si-i2p-plugin \
		-t si-i2p-plugin

fedora:
	docker build --force-rm -f Dockerfile/Dockerfile.build-fedora -t si-i2p-plugin-rpm .
	docker run --name si-i2p-plugin-rpm -t si-i2p-plugin-rpm
	docker exec -t si-i2p-plugin-rpm ls /home/sii2pplugin/
	#docker cp si-i2p-plugin-rpm:/home/sii2pplugin/
	docker rm -f si-i2p-plugin-rpm

checkinstall-rpm: all preinstall-pak postremove-pak description-pak
	checkinstall -R --default \
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
