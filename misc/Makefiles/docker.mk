
docker-clobber:
	docker rm -f si-proxy \
		sam-browser \
		sam-host; true
	docker rmi -f eyedeekay/sam-host \
		eyedeekay/sam-browser \
		eyedeekay/si-i2p-plugin; true
	docker network rm si; true

docker-setup: docker docker-network docker-host docker-run docker-browser

docker-browser:
	docker build --force-rm \
		--build-arg BROWSER_VERSION="$(BROWSER_VERSION)" \
		--build-arg PORT="$(BROWSER_PORT)" \
		-f Dockerfiles/Dockerfile.browser -t eyedeekay/sam-browser .

browse: docker-browser
	docker run --rm -i -t -d \
		-e DISPLAY=$(DISPLAY) \
		-e VERSION="$(BROWSER_VERSION)" \
		--name sam-browser \
		--network si \
		--network-alias sam-browser \
		--hostname sam-browser \
		--link si-proxy \
		--ip 172.80.80.5 \
		--volume /tmp/.X11-unix:/tmp/.X11-unix:ro \
		--volume $(browser):/home/anon/tor-browser_en-US/Browser/Desktop \
		eyedeekay/sam-browser sudo -u anon /home/anon/i2p-browser_en-US/Browser/start-i2p-browser \
		$(browse_args)

docker:
	docker build --force-rm -f Dockerfiles/Dockerfile.samhost -t eyedeekay/sam-host .
	docker build --force-rm -f Dockerfiles/Dockerfile.jumphelper -t eyedeekay/sam-jumphelper .
	docker build --force-rm -f Dockerfile -t eyedeekay/si-i2p-plugin .


docker-network:
	docker network create --subnet 172.80.80.0/29 si; true

docker-host:
	docker run \
		-d \
		--name sam-host \
		--network si \
		--network-alias sam-host \
		--hostname sam-host \
		--link si-proxy \
		--restart always \
		--ip 172.80.80.2 \
		-p :4567 \
		-p 127.0.0.1:$(different_port):7073 \
		--volume $(i2pd_dat):/var/lib/i2pd:rw \
		-t eyedeekay/sam-host; true

docker-jumphelper:
	docker run \
		-d \
		--name sam-jumphelper \
		--network si \
		--network-alias sam-jumphelper \
		--hostname sam-jumphelper \
		--link si-proxy \
		--restart always \
		--ip 172.80.80.3 \
		-p 127.0.0.1:7054:7054 \
		-t eyedeekay/sam-host; true

docker-run: docker-tidy docker-host
	@sleep 1
	docker run \
		-d \
		--name si-proxy \
		--network si \
		--network-alias si-proxy \
		--hostname si-proxy \
		--link sam-host \
		--link sam-browser \
		--link sam-jumphelper \
		--user sii2pplugin \
		--ip 172.80.80.4 \
		-p 127.0.0.1:44443:44443 \
		-p 127.0.0.1:44446:44446 \
		--restart always \
		-t eyedeekay/si-i2p-plugin

docker-follow:
	docker logs -f si-proxy

docker-tidy:
	docker rm -f si-proxy; true

docker-clean: docker-clean
	docker rm -f sam-host; true
	docker rmi -f eyedeekay/si-i2p-plugin; true

docker-copy:
	docker cp sam-browser:/home/anon/i2p-browser.tar.gz ../di-i2p-browser.tar.gz

stop:
	docker rm -f si-proxy; true
