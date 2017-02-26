Stream-Isolating i2p HTTP Proxy(SAM Application)
------------------------------------------------

This is an i2p SAM application which presents an http proxy(on port 4443 by
default) that acts as an intermediate between your browser and the i2p network.
Then it uses the SAM library to create a unique destination for each i2p site
that you visit. This way, your base32 destination couldn't be used to track you
with a network of colluding sites. I doubt it's a substantial problem right now
but it might be someday.

How it will work:
=================

First it sets up an HTTP proxy on your local machine.

        [ HTTP Proxy ]

This HTTP Proxy is used to organize Tunnels, which are paths between i2p
destinations.

        [ HTTP Proxy ]
                      [List Of Tunnels]

This HTTP Proxy intercepts your requests and checks to see if you have already
connected to an eepSite.

        Request->[ HTTP Proxy ]
                         |                        [List Of Tunnels]
                         +->New eepSite
                         |
                         +->Visited eepSite

If you haven't connected to the eepSite before, it creates a new tunnel specific
to that eepSite by contacting the SAM bridge. Once that is done, the request is
sent using the new tunnel.

        Request->[ HTTP Proxy ]
                         |
                         +->New eepSite+
                         |             |
                         |             +[List Of Tunnels + new Tunnel]
                         |                                   |
                         |                                   +->[List Of Tunnels]:Request
                         +->Visited eepSite

If you've already connected to the eepSite, it makes the request using the
destination already associated with the eepSite.

        Request->[ HTTP Proxy ]
                         |
                         +->New eepSite
                         |
                         +->Visited eepSite
                                    |
                                    +->[List Of Tunnels]:Request

Right now it's a work in progress, but it should only take a couple days to do.

Still non-functional, but usage is starting to be defined. So far:

        -bridge="host:port of the SAM bridge(requires both)(defaults to localhost:7656)"
        -proxy="host:port of http proxy(requires both)(defaults to localhost:4443)"
        -log="path to log file(defaults to $HOME/.i2pstreams.log)"

