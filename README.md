Stream-Isolating i2p HTTP Proxy(SAM Application)
------------------------------------------------

This is an i2p SAM application which presents an HTTP proxy(on port 4443 by
default) that acts as an intermediate between your browser and the i2p network.
Then it uses the SAM library to create a unique destination for each i2p site
that you visit. This way, your base32 destination couldn't be used to track you
with a network of colluding sites. I doubt it's a substantial problem right now
but it might be someday. Facebook has an onion site, and i2p should have
destination isolation before there is a facebook.i2p.

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
        -proxy="host:port of HTTP proxy(requires both)(defaults to localhost:4443)"
        -log="path to log file(defaults to $HOME/.i2pstreams.log)"
        -incognito="disables logging, clears destinations"

About the obscure and hypothetical attack this will likely be inadequate to fully prevent
=========================================================================================

A threat to anonymity systems clients is the so-called "Fingerprinting" attack
and, as it turns out, i2p a person using the i2p HTTP proxy to browse the i2p
overlay network can be fingerprinted by the destination of the HTTP proxy.

Wait, that's alot of big words
------------------------------

OK so this speaks to the nuances of discussing anonymity. For the purposes of
this attack, anonymity is any identifying characteristic that can be used to
link your browsing to you. This can be things like I.P. addresses, browser
characteristics, details about your hardware, or details about your running
system. i2p is pretty good at preventing most of this from getting out, but in
one case, it could make the problem worse.

i2p includes an HTTP proxy for clients(Web browsers) to use when browsing the
i2p network. When you start this HTTP proxy it creates a "destination" which is
the cryptographic identifier used by services to reply to your browser's
requests. This feature of your system is linkable to your identity, which isn't
bad in-and-of-itself, but unfortunately the same destination will be replied
to by all the sites you visit. They will all see the same destination. That's
not good.

Now, I suppose the worst-case scenario is that an anonymous journalist or
activist could be linked to his content by the confiscation of his i2p router.
But other things could happen too. If i2p web-services became more popular as a
way to communicate, perhaps with self-hosted social networks, a social-network
operator could set up a network of colluding sites to gather information about
the social-network account associated with a particular router's HTTP proxy
destination. This could allow them to build advertising profiles without the
consent of a user, and potentially link a user to his or her private
communications.

So what does this do about it?
------------------------------

This is a special i2p HTTP proxy that talks to the SAM(Simple Anonymous
Messaging) bridge, an i2p API for developing application-layer software. After
initializing it's connection with the SAM bridge, it it creates a tunnel session
which connects to i2p but doesn't retrieve any information. Then you make a
request, and it checks whether you've visited that site yet or not. If you have
not, it will connect that requested destination to that incomplete tunnel and
generate a new destination for the next site. If you have visited the site, it
will use the existing destination. This gives you the ability to have a unique
identity on a per-site basis and prevents you from being linkable across
multiple eepSites.

What doesn't it do?
===================

Filtering. It doesn't remove unnecessary headers, Javascript, or regularize your
user agent string. It will never filter Javascript on it's own, but it will make
some attempt to filter headers and optionally rewrite the user-agent string(with
the default being the current TBB string)before I consider it ready.
