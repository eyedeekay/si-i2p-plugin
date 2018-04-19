
test-pipes:
	make test-easy; sleep 1
	make test-less; sleep 1
	make test-real; sleep 1
	make test-hard; sleep 1
	make test-diff; sleep 1
	make test-dfhd; sleep 1
	make test-dfsd; sleep 1

test-fake:
	@echo " It should not simply crash upon recieving a bad request, instead"
	@echo "it should log it, not make the request or touch the network at"
	@echo "all, and move on."
	echo http://notarealurl.i2p > parent/send

test-leak:
	@echo " When a client requests a clear web resource, perhaps by accident,"
	@echo "the proxy should immediately refuse to fetch the data and not leak"
	@echo "information to network services."
	echo http://duckduckgo.com > parent/send

test-less:
	@echo " It should not simply crash upon recieving a bad request, instead"
	@echo "it should log it, not make the request or touch the network at"
	@echo "all, and move on. This should include urls that don't exist under"
	@echo "domains that do."
	echo http://i2p-projekt.i2p/en/download > parent/send

test-easy:
	@echo " It should know how to send requests for well-formed http url's"
	@echo "that point to b32 addresses or sites in the address book"
	echo http://i2p-projekt.i2p > parent/send

test-hard:
	@echo " It should also be able to recognize and correct simple"
	@echo "formatting mistakes in URL's and correct them where appropriate."
	echo i2p-projekt.i2p > parent/send

test-real:
	@echo " It should also be able to recognize and correct simple"
	@echo "formatting mistakes in URL's and correct them where appropriate."
	echo i2p-projekt.i2p/en/download > parent/send

test-sf: test-real test-realcurl

test-df: test-diff test-dfhd test-dfsd

test-diff:
	@echo " It should know how to send requests for well-formed http url's"
	@echo "that point to b32 addresses or sites in the address book"
	echo http://inr.i2p > parent/send

test-dfhd:
	@echo " It should know how to send requests for well-formed http url's"
	@echo "that point to b32 addresses or sites in the address book"
	echo inr.i2p > parent/send

test-dfsd:
	@echo " It should know how to send requests for well-formed http url's"
	@echo "that point to b32 addresses or sites in the address book"
	echo inr.i2p/latest > parent/send

test-loop:
	@echo " It's rude and a privacy risk to use i2p-projekt.i2p(or any" "
	@echo "nonlocal eepSite as a test url in the final application. Instead,"
	@echo "it should generate a destination on this machine and query it as"
	@echo "a test, then immediately tear down the test tunnel."
	echo test.i2p > parent/serv
	echo http://test.i2p > parent/send

test-proxy: test-realcurl test-http test-curl test-curldiff test-httpdiff test-httpdiffsd

test-curl:
	@echo "Test the http proxy in as simple a way as possible"
	/usr/bin/curl -w "@misc/curl-format.txt" -x 127.0.0.1:4443 i2p-projekt.i2p

test-curl-css:
	/usr/bin/curl -w "@misc/curl-format.txt" -x 127.0.0.1:4443 i2p-projekt.i2p/_static/styles/duck/syntax.css

test-http:
	@echo "Test the http proxy in as simple a way as possible"
	/usr/bin/curl -w "@misc/curl-format.txt" -x 127.0.0.1:4443 http://i2p-projekt.i2p

test-realcurl:
	/usr/bin/curl -w "@misc/curl-format.txt" -x 127.0.0.1:4443 http://i2p-projekt.i2p/en/download

test-curldiff:
	@echo "Test the http proxy in as simple a way as possible"
	/usr/bin/curl -w "@misc/curl-format.txt" -x 127.0.0.1:4443 inr.i2p

test-httpdiff:
	@echo "Test the http proxy in as simple a way as possible"
	/usr/bin/curl -w "@misc/curl-format.txt" -x 127.0.0.1:4443 http://inr.i2p

test-httpdiffsd:
	@echo "Test the http proxy in as simple a way as possible"
	/usr/bin/curl -w "@misc/curl-format.txt" -x 127.0.0.1:4443 http://inr.i2p/latest

test-proxy-time: clean-time test-realcurl-time test-http-time test-curl-time test-curldiff-time test-httpdiff-time test-httpdiffsd-time

clean-time:
	 rm -f misc/timed-test-newproxy.txt

test-curl-time:
	@echo "timing request to i2p-projekt.i2p/ via curl" | tee -a misc/timed-test-newproxy.txt
	/usr/bin/curl -o /dev/null -w "@misc/curl-format.txt" -x 127.0.0.1:4443 i2p-projekt.i2p | tee -a misc/timed-test-newproxy.txt

test-http-time:
	@echo "timing request to http://i2p-projekt.i2p/ via curl" | tee -a misc/timed-test-newproxy.txt
	/usr/bin/curl -o /dev/null -w "@misc/curl-format.txt" -x 127.0.0.1:4443 http://i2p-projekt.i2p | tee -a misc/timed-test-newproxy.txt

test-realcurl-time:
	@echo "timing request to http://i2p-projekt.i2p/en/download via curl" | tee -a misc/timed-test-newproxy.txt
	/usr/bin/curl -o /dev/null -w "@misc/curl-format.txt" -x 127.0.0.1:4443 http://i2p-projekt.i2p/en/download | tee -a misc/timed-test-newproxy.txt

test-curldiff-time:
	@echo "timing request to inr.i2p/ via curl" | tee -a misc/timed-test-newproxy.txt
	/usr/bin/curl -o /dev/null -w "@misc/curl-format.txt" -x 127.0.0.1:4443 inr.i2p | tee -a misc/timed-test-newproxy.txt

test-httpdiff-time:
	@echo "timing request to http://inr.i2p/ via curl" | tee -a misc/timed-test-newproxy.txt
	/usr/bin/curl -o /dev/null -w "@misc/curl-format.txt" -x 127.0.0.1:4443 http://inr.i2p | tee -a misc/timed-test-newproxy.txt

test-httpdiffsd-time:
	@echo "timing request to http://inr.i2p/latest via curl" | tee -a misc/timed-test-newproxy.txt
	/usr/bin/curl -o /dev/null -w "@misc/curl-format.txt" -x 127.0.0.1:4443 http://inr.i2p/latest | tee -a misc/timed-test-newproxy.txt

test-browser:
	http_proxy=http://127.0.0.1:4443 surf http://inr.i2p

test-browser-diff:
	http_proxy=http://127.0.0.1:4443 surf http://i2p-projekt.i2p

thirdeye:
	http_proxy=http://127.0.0.1:4443 surf http://lxik2bjgdl7462opwmkzkxsx5gvvptjbtl35rawytkndf2z7okqq.b32.i2p/index.html

jump:
	http_proxy=http://127.0.0.1:4443 surf http://lxik2bjgdl7462opwmkzkxsx5gvvptjbtl35rawytkndf2z7okqq.b32.i2p/jump/i2pforum.i2p

curljump:
	/usr/bin/curl -X 127.0.0.1:4443 lxik2bjgdl7462opwmkzkxsx5gvvptjbtl35rawytkndf2z7okqq.b32.i2p/jump/i2pforum.i2p

surf:
	http_proxy=http://127.0.0.1:4443 surf http://i2p-projekt.i2p

firefox:
	iceweasel http://i2p-projekt.i2p > firefox.log 2> firefox.err &
