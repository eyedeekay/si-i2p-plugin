
demoservice:
	docker build -f Dockerfiles/Dockerfile.demoservice -t eyedeekay/i2p-demoservice .

democlean:
	docker rm -f demoservice; true
	rm -f demo.b32.i2p

demo: democlean demoservice
	docker run -d --cap-drop all --name demoservice -p :4567 -p 7071:7070  -p 7656:7658 -t eyedeekay/i2p-demoservice
	sleep 10
	make democonfig
	@echo "WARNING: going to sleep for 20 minutes to allow new eepSite to become available"
	sleep 20m

democonfig: demo.b32.i2pclean demo.b32.i2p

demo.b32.i2pclean:
	rm -f demo.b32.i2p

demo.b32.i2p:
	/usr/bin/lynx -dump -listonly http://localhost:7071/?page=i2p_tunnels | grep b32 | sed 's| 12||g' | sed 's| 11||g' | sed 's| 10||g' | sed 's| 9||g' | sed 's| 8||g' | sed 's|http://localhost:7071/?page=local_destination&b32=||g' |  tr -d ' .' | tee demo.b32.i2p

demo-1-vuln:
	/usr/bin/curl -x 127.0.0.1:4444 $(shell head -n 1 demo.b32.i2p).b32.i2p

demo-2-vuln:
	/usr/bin/curl -x 127.0.0.1:4444 $(shell tail -n 1 demo.b32.i2p).b32.i2p

demo-1-fix:
	/usr/bin/curl -x 127.0.0.1:4443 $(shell head -n 1 demo.b32.i2p).b32.i2p

demo-2-fix:
	/usr/bin/curl -x 127.0.0.1:4443 $(shell tail -n 1 demo.b32.i2p).b32.i2p

headers.vuln.log: misc/log/headers.vuln.log

misc/log/headers.vuln.log:
	mkdir -p misc/log
	docker logs demoservice | tee misc/log/headers.vuln.log

demo-p1-fix:
	@echo $(shell head -n 1 demo.b32.i2p).b32.i2p | tee parent/send

demo-p2-fix:
	@echo $(shell tail -n 1 demo.b32.i2p).b32.i2p | tee parent/send

demo-pipes: democonfig demo-p1-fix demo-p2-fix headers.fix.log

headers.fix.log: misc/log/headers.fix.log

misc/log/headers.fix.log:
	mkdir -p misc/log
	docker logs demoservice | tee misc/log/headers.fix.log

demo-vuln: democonfig demo-1-vuln demo-2-vuln headers.vuln.log
	@echo "Un-Isolated demo completed: Note the X-I2P-DEST* headers are the same between sites"

demo-fix: democonfig demo-1-fix demo-2-fix headers.fix.log
	@echo "Un-Isolated demo completed: Note the X-I2P-DEST* are different now"

demo-all: demo-vuln demo-fix

demolog:
	rm -rf misc/log
	mkdir -p misc/log

dodemo: demolog run demo
	make demo-all
