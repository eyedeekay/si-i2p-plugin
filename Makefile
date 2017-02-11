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
	./si-i2p-plugin
