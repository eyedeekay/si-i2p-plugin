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
