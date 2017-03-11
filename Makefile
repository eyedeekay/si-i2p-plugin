PREFIX := /usr/local/bin

all:
	go build

install:
	cp si-i2p-plugin ${PREFIX}/si-i2p-plugin

clean:
	rm si-i2p-plugin

uninstall:
	rm ${$PREFIX}/si-i2p-plugin

tryout:
	make
	gdb si-i2p-plugin

memcheck:
	make
	valgrind --track-origins=yes --leak-check=full ./si-i2p-plugin

lint:
	golint

testreq:
	wget -qO - -e use_proxy=yes -e http_proxy=127.0.0.1:4444 http://i2p-projekt.i2p

request:
	wget -qO - -e use_proxy=yes -e http_proxy=127.0.0.1:4443 http://i2p-projekt.i2p

req:
	make testreq
	make request

cmds:
	cat Makefile
