Stream-Isolating i2p HTTP Proxy(SAM Application)
------------------------------------------------

This is an i2p SAM application which presents an http proxy(on port 4443 by
default) that acts as an intermediate between your browser and the i2p network.
Then it uses the SAM library to create a unique destination for each i2p site
that you visit. This way, your base32 destination couldn't be used to track you
with a network of colluding sites. I doubt it's a substantial problem right now
but it might be someday.

Right now it's a work in progress, but it should only take a couple days to do.