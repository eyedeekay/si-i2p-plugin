

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
	install init.d/si-i2p-plugin $(PREFIX)$(ETC)init.d/si-i2p-plugin
	install systemd/sii2pplugin.service $(PREFIX)$(ETC)systemd/system
	install si-i2p-plugin/settings.cfg $(PREFIX)$(ETC)si-i2p-plugin/settings.cfg

try:
	bash -c "./bin/si-i2p-plugin 1>log 2>err" & sleep 1 && true
	cat i2p-projekt.i2p/recv | tee test.html

test:
	echo http://i2p-projekt.i2p > i2p-projekt.i2p/send

alttest:
	./bin/si-i2p-plugin --url fireaxe.i2p

testother:
	echo http://i2p-projekt.i2p/en/download > i2p-projekt.i2p/send

clean:
	rm -rf i2p-projekt.i2p bin/si-i2p-plugin *.html *-pak *err *log

cat:
	cat i2p-projekt.i2p/recv

name:
	cat i2p-projekt.i2p/name

exit:
	echo y > i2p-projekt.i2p/del

noexit:
	echo n > i2p-projekt.i2p/del

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

checkinstall:
	make preinstall-pak
	make postremove-pak
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
	mkdir -p preinstall-pak
	@echo "#! /bin/sh" | tee preinstall-pak/preinstall
	@echo "adduser --system --no-create-home --disabled-password --disabled-login --group sii2pplugin" | tee -a preinstall-pak/preinstall
	@echo "chown -R sii2pplugin:adm $(PREFIX)$(VAR)$(LOG)/si-i2p-plugin/ $(PREFIX)$(VAR)$(RUN)si-i2p-plugin/" | tee -a preinstall-pak/preinstall

postremove-pak:
	mkdir -p postremove-pak
	@echo "#! /bin/sh" | tee postremove-pak/postremove
	@echo "deluser sii2pplugin" | tee -a postremove-pak/postremove

docker:
	docker build -f Dockerfile/Dockerfile -t si-i2p-plugin .

fedora:

