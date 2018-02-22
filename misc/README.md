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
