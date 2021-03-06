Destination-Isolating i2p HTTP Proxy(SAM Application)
=====================================================

*one eepSite, one destination. Semi-automatic contextual identity management*
*for casually browsing i2p eepSites.*

This is an i2p SAM application which presents an HTTP proxy(on port 4443 by
default) that acts as an intermediate between your browser and the i2p network.
Then it uses the SAM library to create a unique destination for each i2p site
that you visit. This way, your unique destination couldn't be used to track you
with a network of colluding sites. I doubt it's a substantial problem right now
but it might be someday. Facebook has an onion site, and i2p should have
destination isolation before there is a facebook.i2p.

[I2P Link - Stream Isolation](http://trac.i2p2.i2p/ticket/1149)
[I2P Link - Shared Tunnels](http://zzz.i2p/topics/217)

Development here has ceased. A replacement that will work much better is in
development at [eyedeekay/eeProxy](https://github.com/eyedeekay/eeProxy).

[![Build Status](https://travis-ci.org/eyedeekay/si-i2p-plugin.svg?branch=master)](https://travis-ci.org/eyedeekay/si-i2p-plugin)

[Usage](USAGE.md)
-----------------

License
-------

Sorry for forgetting to license this. It's MIT Licensed now.

Installation:
-------------

Moved to [misc/docs/INSTALL.md](misc/docs/INSTALL.md)

### Usage Examples(HTTP Proxy)

#### firefox

![Firefox Configuration](misc/firefox.png)

#### curl

        curl -x 127.0.0.1:4443 http://i2p-projekt.i2p

#### surf

        export http_proxy="http://127.0.0.1:4443" surf http://i2p-projekt.i2p

What works so far:
------------------

### It seems to do exactly what it says on the package.

If you'd like to test it, the easiest way is to use Docker. To generate all
the required containers locally and start a pre-configured browser, run:

        make docker-setup browse

### The http proxy

It's *still a little experimental*, but currently it is possible to set
your web browser's HTTP proxy to localhost:4443 and use it to browse eepSites
with the guarantee that you will show a different destination to every eepSite
you visit. Combined with a browser designed to minimize your uniqueness to the
servers you visit, you should be much more difficult to track across sites and
across sessions.

**That said there *is the obvious thing to consider*:**

If it wasn't super, super obvious to everyone, it's really, really easy to tell
the difference between this proxy and the default i2p/i2pd http proxies and I
don't think there's anything I can do about that. Also *if you're the only*
*person to visit a particular group of colluding eepSites* then it's *still*
*possible to link your activities by timing*, but *I don't think it's possible*
to "*prove*" that it the same person exactly(certainly not in a cryptographic
sense), just that it's likely to be the same person. I can't do anything about
a small anonymity set. That said, the tunnels created by this proxy are
inherently short-lived and tear themselves down after an inactivity timeout,
requiring an attacker to request resources over and over to keep a tunnel alive
long-term to be useful for tracking. By the way, as far as I know, using this
will drastically reduce your anonymity set unless it's widely adopted. TESTING
ONLY.

I haven't been able to crash it or attack it by adapting known attacks on
browsers and HTTP proxies to this environment. It should at least fail early if
something bad happens.

#### User-Defined Jump Hosts

Addresshelpers are tentatively working again. I'll need to explain more on how
soon.

#### Current Concerns:

It mostly looks like Jumphelper is working immediately after it's started, but
if you try to visit a base32 address before they sync, then si-i2p-plugin will
ask jumphelper for it and fail, interpreting it as a hostname error. It lasts
for like, 15 seconds and only affects base32 URLs and is fixable.

[I wonder if I could make it talk to TorButton?](https://www.torproject.org/docs/torbutton/en/design/index.html.en)

Littleboss seems like it could potentially be useful in lots of contexts, like
in a portable browser bundle and stuff like that.

I really need to get a Windows machine, but it's just such an awful way to work.

What the final version should do:
---------------------------------

The final version should use the parent pipe and the aggregating pipe to send
and recieve requests as an http proxy in the familiar way.

Version Roadmap:

  * ~~0.17 - Named Pipes work for top-level i2p domains and can retrieve~~
   ~~directories under a site~~
  * ~~0.18 - Named Pipes for i2p domains and can retrieve subdirectories,~~
   ~~which it caches in clearly-named folders as normal files(Containing HTML)~~
  * ~~0.19 - Expose an http proxy that hooks up to the existing infrastructure~~
   ~~for destination isolation~~
  * 0.20 - ~~Ready for more mainstream testing~~, ~~should successfully isolate~~
   ~~requests for resources embedded in the retrieved web pages. Addresshelper~~
   ~~functions but it's it's own whole thing now.~~
  * 0.21 - Should be able to generate services on the fly by talking to the SAM
   bridge. Truth be told, except when I break something, it works nearly
   perfectly. I do break stuff alot though, development is quite active and I
   have to leverage the communication tools available to me to get it done, of
   which github is, for better or for worse, an integral one at this point.
   ~~Tunnels should kill themselves after inactivity, and revive themselves with~~
   ~~new identities after. This will help minimize the impact of cross-site~~
   ~~resource-request timing attacks by making destinations more ephemeral,~~
   ~~requiring an attacker to keep tunnels alive to monitor an identity~~
   ~~long-term.~~
  * 0.22 - Library-fication should be finished by here. Turning the underlying
   code into a library will mostly be a matter of identifying which features need
   to be exposed for it to be useful in that way. I'll update the number when
   I've written go-based tests for it. ~~Maybe 3/4ths of the code that exists~~
   ~~has relevant tests now. All of them actually test a thing that knows how~~
   ~~to fail now.~~
  * 0.23 - ~~Enable additional configuration options, like tunnel lengths~~
   ~~(always symmetrical) tunnel quantities(not always symmetical) idle~~
   ~~connections per host, and backup tunnel quantity.~~ Implement "Saved
   Tunnels" which will allow per-eepSite long-term destinations, and ironically,
   potential parity with standard http proxies from i2p and i2pd. I used that
   word deliberately.
  * 0.24 - Experiment with adding a SOCKS proxy. Create a version which contains
   a SOCKS proxy for testing. Actually have a SOCKS proxy. [This should be acceptable in implementing the SOCKS proxy](https://github.com/armon/go-socks5)
   Torbutton Control Port compatibility.
  * 0.25 - Package.
  * 0.26 - Android? WebRTC over SSU? Try linking libi2pd? Help make go-i2p into
   a usable router? Browser Bundle(s?)? At risk of growing a little further than
   I originally intended, but I feel like there's some potential here.

### The pipes

Moved to [misc/docs/PIPES.md](misc/docs/PIPES.md)

Screenshots:
------------

[moved to misc/SCREENSHOTS.md](misc/SCREENSHOTS.md)

Donate
------

### Monero Wallet Address

  XMR:43V6cTZrUfAb9JD6Dmn3vjdT9XxLbiE27D1kaoehb359ACaHs8191mR4RsJH7hGjRTiAoSwFQAVdsCBToXXPAqTMDdP2bZB

### Bitcoin Wallet Address

  BTC:159M8MEUwhTzE9RXmcZxtigKaEjgfwRbHt
