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

[**i2p link** A reference to this issue found on zzz.i2p, that I should have put in the readme sooner.](http://zzz.i2p/topics/217)

What works so far:
------------------

### It seems to do exactly what it says on the package.

If you'd like to test it, the easiest way is to use Docker. To generate all
the required containers locally and start a pre-configured browser, run:

        make docker-setup browse

### The http proxy

Again, *still pretty experimental*, but currently it is possible to set
your web browser's HTTP proxy to localhost:4443 and use it to browse eepSites.
I haven't been able to crash it or attack it by adapting known attacks on
browsers and HTTP proxies to this environment. It should at least fail early if
something bad happens.

#### User-Defined Jump Hosts

Addresshelper/Jump hosts are broken in here. I'm working on it.

#### Examples

##### firefox

![Firefox Configuration](misc/firefox.png)

##### curl

        curl -x 127.0.0.1:4443 http://i2p-projekt.i2p

##### surf

        export http_proxy="http://127.0.0.1:4443" surf http://i2p-projekt.i2p

#### Current Concerns:

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

I am now fairly certain that it can't be forced to retrieve URL's outside the
i2p network in properly configured browsers under normal circumstances. Remember
to set [*] Use this proxy server for all protocols or other relevant browser
configurations. This appears to be the default behavior for surf and uzbl.

Before version 0.21, a framework for generating service tunnels ad-hoc will also
be in place. This will be used for fuzz-testing the http proxy and the pipe
proxy. Almost everything will be improved by the availability of this. Before
version 0.25, whether to use either in or out pipes or to enable the pipes at
all, will be configurable.

[I wonder if I could make it talk to TorButton?](https://www.torproject.org/docs/torbutton/en/design/index.html.en)

Elephant in the room #1, it's kind of unfortunately named. I really have a knack
for that.

Elephant in the room #2, it runs excellent on anything that can work with the
named pipe implementation in regular Go. I could take shortcuts that would limit
the functionality available to Windows people, or figure out some way to
implement that functionality on a per-platform basis without losing
functionality. Oh shit conditional compilation in go is super easy! An early
Windows version is available, but everything that's a named pipe in a Unix is a
real file in Windows. So only use the HTTP proxy. Ever. At least until I find a
way to ensure that sent requests are cleared from the file. Preliminary Windows
support is enabled by turning the FIFO's into files and specifying their
behavior in a windows-only version of si-fs-helpers.go. If this turns out to be
good enough then this is how I'll keep doing it.

Elephant in the room #3, absolutely ZERO outproxy support. But that's not really
what it's for. It's probably 90% unhelpful for outproxies anyway.

### The pipes

Moved to [misc/docs/PIPES.md](misc/docs/PIPES.md)

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
   ~~requests for resources embedded in the retrieved web pages.~~ Addresshelper
   still needs alot of work. I did this too soon.
  * 0.21 - Should be able to generate services on the fly by talking to the SAM
  bridge. First worthwhile release for people who aren't shell enthusiasts.
  ~~Tunnels should kill themselves after inactivity, and revive themselves with~~
  ~~new identities after. This will help minimize the impact of cross-site~~
  ~~resource-request timing attacks by making destinations more ephemeral,~~
  ~~requiring an attacker to keep tunnels alive to monitor an identity~~
  ~~long-term.~~
  * 0.22 - Library-fication should be finished by here. Turning the underlying
  code into a library will mostly be a matter of identifying which features need
  to be exposed for it to be useful in that way. I'll update the number when
  I've written go-based tests for it. ~~Maybe 1/5th of it has relevant tests~~
  ~~now.~~
  * 0.23 - ~~Enable additional configuration options, like tunnel lengths~~
  ~~(always symmetrical) tunnel quantities(not always symmetical) idle~~
  ~~connections per host, and backup tunnel quantity.~~ If I'm being honest,
  this will probably be done before 0.21 and 0.22, but it won't be incremented
  until they are done too.
  * 0.24 - Experiment with adding a SOCKS proxy. Create a version which contains
  a SOCKS proxy for testing. Actually have a SOCKS proxy. [This should be acceptable in implementing the SOCKS proxy](https://github.com/armon/go-socks5)
  Torbutton Control Port compatibility.
  * 0.25 - Package.
  * 0.26 -


Installation and Usage:
=======================

Moved to [misc/docs/INSTALL.md](misc/docs/INSTALL.md)

Screenshots:
------------

[moved here:](misc/SCREENSHOTS.md)

Donate
------

### Monero Wallet Address

  XMR:43V6cTZrUfAb9JD6Dmn3vjdT9XxLbiE27D1kaoehb359ACaHs8191mR4RsJH7hGjRTiAoSwFQAVdsCBToXXPAqTMDdP2bZB
