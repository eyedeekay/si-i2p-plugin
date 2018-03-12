
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

static-end-pipe:
	@echo gkso46tc47hdua2kva5zahj3unmyh6ia7bv5oc2ybn2hmeowpz7a.b32.i2p | tee parent/send

headers.fix.log: misc/log/headers.fix.log

misc/log/headers.fix.log:
	mkdir -p misc/log
	docker logs demoservice | tee misc/log/headers.fix.log

demo-pipes: demo-p1-fix demo-p2-fix demo-p3-fix

demo-all: demo-vuln demo-fix

demolog:
	rm -rf misc/log
	mkdir -p misc/log

test-proxy-time-old: clean-time-old test-realcurl-time-old test-http-time-old test-curl-time-old test-curldiff-time-old test-httpdiff-time-old test-httpdiffsd-time-old

clean-time-old:
	 rm -f misc/timed-test-oldproxy.txt

test-curl-time-old:
	@echo "timing request to i2p-projekt.i2p/ via curl" | tee -a misc/timed-test-oldproxy.txt
	/usr/bin/curl -o /dev/null -w "@misc/curl-format.txt" -x 127.0.0.1:4444 i2p-projekt.i2p | tee -a misc/timed-test-oldproxy.txt

test-http-time-old:
	@echo "timing request to http://i2p-projekt.i2p/ via curl" | tee -a misc/timed-test-oldproxy.txt
	/usr/bin/curl -o /dev/null -w "@misc/curl-format.txt" -x 127.0.0.1:4444 http://i2p-projekt.i2p | tee -a misc/timed-test-oldproxy.txt

test-realcurl-time-old:
	@echo "timing request to http://i2p-projekt.i2p/en/download via curl" | tee -a misc/timed-test-oldproxy.txt
	/usr/bin/curl -o /dev/null -w "@misc/curl-format.txt" -x 127.0.0.1:4444 http://i2p-projekt.i2p/en/download | tee -a misc/timed-test-oldproxy.txt

test-curldiff-time-old:
	@echo "timing request to inr.i2p/ via curl" | tee -a misc/timed-test-oldproxy.txt
	/usr/bin/curl -o /dev/null -w "@misc/curl-format.txt" -x 127.0.0.1:4444 inr.i2p | tee -a misc/timed-test-oldproxy.txt

test-httpdiff-time-old:
	@echo "timing request to http://inr.i2p/ via curl" | tee -a misc/timed-test-oldproxy.txt
	/usr/bin/curl -o /dev/null -w "@misc/curl-format.txt" -x 127.0.0.1:4444 http://inr.i2p | tee -a misc/timed-test-oldproxy.txt

test-httpdiffsd-time-old:
	@echo "timing request to http://inr.i2p/latest via curl" | tee -a misc/timed-test-oldproxy.txt
	/usr/bin/curl -o /dev/null -w "@misc/curl-format.txt" -x 127.0.0.1:4444 http://inr.i2p/latest | tee -a misc/timed-test-oldproxy.txt
