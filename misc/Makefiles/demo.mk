
demoservice:
	docker build -f Dockerfiles/Dockerfile.demoservice -t eyedeekay/i2p-demoservice .

democlean:
	docker rm -f demoservice; true
	rm -f demo.b32.i2p

demo: democlean demoservice
	docker run -d --cap-drop all --name demoservice -p :4567 -p 7071:7070 -t eyedeekay/i2p-demoservice
	sleep 10
	make demo.b32.i2p
	@echo "WARNING: going to sleep for 20 minutes to allow new eepSite to become available"
	sleep 20m

demo.b32.i2p:
	/usr/bin/lynx -dump -listonly http://localhost:7071/?page=i2p_tunnels | grep b32 | sed 's| 10||g' | sed 's| 9||g' | sed 's|http://localhost:7071/?page=local_destination&b32=||g' |  tr -d ' .' | tee demo.b32.i2p

demo-1-vuln:
	curl -x 127.0.0.1:4444 $(shell head -n 1 demo.b32.i2p).b32.i2p

demo-2-vuln:
	curl -x 127.0.0.1:4444 $(shell tail -n 1 demo.b32.i2p).b32.i2p

misc/log/headers.vuln.log:
	mkdir -p misc/log
	docker logs demoservice | tee misc/log/headers.vuln.log

demo-1-fix:
	curl -x 127.0.0.1:4443 $(shell head -n 1 demo.b32.i2p).b32.i2p

demo-2-fix:
	curl -x 127.0.0.1:4443 $(shell tail -n 1 demo.b32.i2p).b32.i2p

misc/log/headers.fix.log:
	mkdir -p misc/log
	docker logs demoservice | tee misc/log/headers.fix.log

demo-vuln: demo demo-1-vuln demo-2-vuln misc/log/headers.vuln.log
	@echo "Un-Isolated demo completed: Note the X-I2P-DEST* headers are the same between sites"

demo-fix: demo demo-1-fix demo-2-fix misc/log/headers.fix.log
	@echo "Un-Isolated demo completed: Note the X-I2P-DEST* are different now"

dodemo: run
	make demo-vuln
	make demo-fix
