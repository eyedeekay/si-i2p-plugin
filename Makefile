
build:
	go build -o si-i2p-plugin

install:
	install si-i2p-plugin /usr/local/bin/

try:
	make build
	bash -c "./si-i2p-plugin 1>log 2>err" & sleep 1 && true
	cat i2p-projekt.i2p/name

test:
	echo http://i2p-projekt.i2p/en/docs./api/samv3 > i2p-projekt.i2p/send

clean:
	rm -rf i2p-projekt.i2p

cat:
	cat i2p-projekt.i2p/recv

exit:
	echo yes > i2p-projekt.i2p/del
