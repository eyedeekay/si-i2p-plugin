Destination-Isolating i2p HTTP Proxy(SAM Application)
=====================================================

This is an i2p SAM application which presents an HTTP proxy(on port 4443 by
default) that acts as an intermediate between your browser and the i2p network.
Then it uses the SAM library to create a unique destination for each i2p site
that you visit. This way, your base32 destination couldn't be used to track you
with a network of colluding sites. I doubt it's a substantial problem right now
but it might be someday. Facebook has an onion site, and i2p should have
destination isolation before there is a facebook.i2p.

What works so far:
------------------

In the front, right now there are three "Parent" pipes which are used to
delegate requests and order responses from the system which exists behind them
and signal the interruption of the isolating proxy. It can't be hooked up to a
web browser yet, but you might be able to work something out with like, socat or
something. If you run the application ./si-i2p-plugin from this repository it
will create a folder with the name "parent" containing the following named
pipes.

        parent/
                send
                     <- echo "desired-url" > parent/send
                recv
                     <- cat parent/recv
                del
                     <- echo "y" > parent/del

At this point, no connection to either the SAM bridge or the i2p network has
actually been made yet. The parent pipes are simply ready to make the connection
when necessary. In order to make a request, pipe a URL into the parent/send
pipe. To read out the most recent response, cat out the parent/recv pipe. Lastly
to close all the pipes and clean up, echo "y" into parent/del.

Behind that, there is a system which uses named pipes to allow a user to send
and recieve requests and get information about eepSites on the i2p network. If
you were to, for instance, make a request for i2p-projekt.i2p through
parent/send, it would look for the SAM session associated with that site(or
create one if it doesn't exist) in a folder called "i2p-projekt.i2p". Inside
that folder will be 4 files corresponding to the named pipes:

        destination_url.i2p/
                            send
                                 <- echo "desired-url" > destination_url.i2p/send
                            recv
                                 <- cat destination_url.i2p/recv
                            name
                                 <- cat destination_url.i2p/name
                            del
                                 <- echo "y" > destination_url.i2p/del

In order to use them to interact with eepSites, you may either make your
requests to the parent pipes  which will delegate the responses to the child
pipes automatically, or you may manually pipe the destination URL into
destination\_url.i2p/send, and pipe out the result from
destination\_url.i2p/recv. To retrieve the full cryptographic identifier of the
eepSite, pipe out the destination from destination\_url.i2p/name and to close
the pipe, pipe anything at all into destination\_url.i2p/del.

What I'm doing right now:
-------------------------

I'm doing the last of the parent pipe system. Making sure it always signals the
correct child pipe and only the correct child pipe. Whenever I get tired I
try and figure out how to integrate it with an initsystem or a package manager.

What the final version should do:
---------------------------------

The final version should use the parent pipe and the aggregating pipe to send
and recieve requests as an http proxy in the familiar way.

Silly Questions I'm asking myself about how I want it to work:
--------------------------------------------------------------

Should it be possible to disable the http proxy and interact with it only using
named pipes? For right now, I'm leaning toward yes as this could be useful for
ssh clients and similar applications.

Should it do filtering? I really don't think so but if there's a simple way to
strip outgoing information then maybe. I dislike complexity. It's why this has
taken so long.

Wierdish applications consequent to the design:
-----------------------------------------------

Tox-i2p Bridging: p2p networking and Unix enthusiasts will probably recognize
that the ideas in this come from the Suckless Tox client, [Pranomostro's ratox](https://github.com/pranomostro/ratox).
And it's easy to throw strings back and forth between named pipes. Perhaps with
some tooling this could be used to build a bridge between the Tox and i2p
networks without needing to modify the client.

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
