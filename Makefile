
all:
		go build -o ../bin/si-i2p-plugin ./src

install:
	mkdir -p /var/log/si-i2p-plugin/ /var/si-i2p-plugin/ /etc/si-i2p-plugin/
	chown -R sii2pplugin:adm /var/log/si-i2p-plugin/ /var/si-i2p-plugin/
	install bin/si-i2p-plugin /usr/local/bin/
	install bin/si-i2p-plugin.sh /usr/local/bin/
	install init.d/si-i2p-plugin /etc/init.d/si-i2p-plugin
	install si-i2p-plugin/settings.cfg /etc/si-i2p-plugin/settings.cfg

try:
	bash -c "./bin/si-i2p-plugin 1>log 2>err" & sleep 1 && true
	cat i2p-projekt.i2p/name

test:
	echo http://i2p-projekt.i2p/en/docs./api/samv3 > i2p-projekt.i2p/send

clean:
	rm -rf i2p-projekt.i2p err log

cat:
	cat i2p-projekt.i2p/recv

name:
	cat i2p-projekt.i2p/name

exit:
	echo y > i2p-projekt.i2p/del

html:
	cat i2p-projekt.i2p/recv > test.html; true

html-test:
	sr W ./test.html

user:
	sudo adduser --system --no-create-home --disabled-password --disabled-login --group sii2pplugin

