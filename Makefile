
PREFIX := /usr/local

all:
	go build -compiler gccgo

install:
	cp si-i2p-plugin $(DESTDIR)${PREFIX}/si-i2p-plugin

clean:
	make
	rm si-i2p-plugin
	rm -rf doc-pak err log *.deb wget-log description-pak *.tgz

uninstall:
	rm $(DESTDIR)${PREFIX}/si-i2p-plugin

tryout:
	make
	gdb si-i2p-plugin

memcheck:
	make
	valgrind --track-origins=yes --leak-check=full ./si-i2p-plugin 2> err

lint:
	golint

testreq:
	wget -qO - -e use_proxy=yes -e http_proxy=127.0.0.1:4444 http://zzz.i2p

lreq:
	wget -qO - -e use_proxy=yes -e http_proxy=127.0.0.1:9999 http://zzz.i2p

request:
	wget -qO - -e use_proxy=yes -e http_proxy=127.0.0.1:4443 http://zzz.i2p

req:
	make testreq
	make request

cmds:
	cat Makefile

signreadme:
	rm README.md.asc
	gpg --clear-sign -u C0CEEE297B5FE45FF610AAC6F05F85FA446C042B README.md

testdeb:
	make
	checkinstall --install=no \
		--pkgversion=0.1 \
		--pkgname=di-i2p-plugin \
		--maintainer="problemsolver@openmailbox.org" \
		--pakdir="../"
