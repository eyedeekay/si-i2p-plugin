Config Files and Demo Scripts:
==============================

This folder contains some configuration files and supplemental information about
the i2p destination-isolating proxy.

Makefiles:
----------

I have a bad habit of putting things in makefiles that I shouldn't. In this
folder I at least attempt to separate them into smaller, more self-explanatory
chunks.

Demo config files:
------------------

This is mostly to contain the stuff required to run the demo. The demo is a
docker container running i2pd with two http services, pointing to the same
python-based web service, which simply logs the headers of every request made
to it. Running

        make dodemo

in the project root directory will generate and run the container, go to sleep
for a while to allow the network to get ready and publish the services, and then
make some requests to the servers over the default http proxy and the demo http
proxy, and once they complete, copy the corresponding logs to misc/logs/.

The echo server is borrowed from [1kastner's github gist](https://gist.github.com/1kastner/e083f9e813c0464e6a2ec8910553e632)
and is unaltered so far. Eventually I'll make it print the destination on the
page.

Attacking the Proxy:
--------------------

### Software Fingerprinting:

This proxy is trivially distinguishable from the default i2p http proxies in at
least 2 ways, probably three, maybe more. I doubt I'm going to be able to do
anything about *all* of them, but the same colluding adversary can identify
*this proxy* by identifying which resources you *don't* request from their site.
In other words, if a site administrator knowingly requests off-site resources,
it can identify destinations created by this proxy by identifying destinations
that *never* request those third party resources.

### Timing Attacks:

See issue #2
