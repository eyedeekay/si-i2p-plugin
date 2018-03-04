
headers.vuln.log: misc/log/headers.vuln.log

misc/log/headers.vuln.log:
	mkdir -p misc/log
	docker logs demoservice | tee misc/log/headers.vuln.log

demo-p1-fix:
	@echo lqnwvwsgio6k53zq6d7r5bpaxuslc45vgsiqo6i3ebshkqpgrnma.b32.i2p | tee parent/send

demo-p2-fix:
	@echo zcofypupen75rdv5zihviweyw5emk2l34idq423kbhj7n3owoe5a.b32.i2p | tee parent/send

demo-p3-fix:
	@echo zjjjd756aucwz3pa2fl4mb3po2wtf752aefpod4gvedwreeox52q.b32.i2p | tee parent/send

headers.fix.log: misc/log/headers.fix.log

misc/log/headers.fix.log:
	mkdir -p misc/log
	docker logs demoservice | tee misc/log/headers.fix.log

demo-pipes: demo-p1-fix demo-p2-fix demo-p3-fix

demo-all: demo-vuln demo-fix

demolog:
	rm -rf misc/log
	mkdir -p misc/log

