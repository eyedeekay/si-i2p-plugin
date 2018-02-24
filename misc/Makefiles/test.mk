
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
	/usr/bin/curl -x 127.0.0.1:4443 i2p-projekt.i2p

test-http:
	@echo "Test the http proxy in as simple a way as possible"
	/usr/bin/curl -x 127.0.0.1:4443 http://i2p-projekt.i2p

test-realcurl:
	/usr/bin/curl -x 127.0.0.1:4443 http://i2p-projekt.i2p/en/download

test-curldiff:
	@echo "Test the http proxy in as simple a way as possible"
	/usr/bin/curl -x 127.0.0.1:4443 inr.i2p

test-httpdiff:
	@echo "Test the http proxy in as simple a way as possible"
	/usr/bin/curl -x 127.0.0.1:4443 http://inr.i2p

test-httpdiffsd:
	@echo "Test the http proxy in as simple a way as possible"
	/usr/bin/curl -x 127.0.0.1:4443 http://inr.i2p/latest
