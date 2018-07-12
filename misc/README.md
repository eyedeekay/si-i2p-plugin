Why I think this is an ok idea
==============================

Attacking the user behind the Default Proxy
-------------------------------------------

The default i2p http proxy does exactly what it's supposed to, in a properly
configured browser. The shortest path to that is to customize a TBB, as far as
I know right now. Bromite or similar may be an interesting avenue to explore in
the future. I think the way it works is reasonable. I'll get to why I think so
a little further on, and why I see my approach as complementary, and not
competing.

In the default i2p http proxy, all applications using the proxy share the same
destination, which is visible to all services that they communicate with. That
means that a few kinds of tracking are possible.

  1. An attacker who controls multiple eepSites can link a user's behavior on
  one eepSite to their behavior on another eepSite if they visit them at the
  same time. [I wrote a turnkey implementation of it here](https://github.com/eyedeekay/colluding-sites-attack)
  2. An attacker who provides a service as an eepSite can presumably link a
  visitor to an identifier(user account, etc). This means that they can link
  old destination behavior on any eepSite they control to old destination
  behavior on any other eepSite they control using a single identifier, on a
  single eepSite, to link behavior across old destinations, provided that
  account is used for that session.
  3. If a user leaves a session on an eepSite open, and that eepSite is able to
  send javascript to be executed by the user's browser, then the eepSite could
  keep the default http proxy from timing out and generating a new destination.
  This could be abused to make a destination exist longer, perhaps making
  tracking easier.

I don't think most of us trust all the eepSites out there *that* much. Some,
sure, but not all of them. If the eepSite ecosystem were to grow, I think the
prevalence of this kind of tracking would probably grow too. Imagine attack #2
if there were a facebook.i2p? I mean they have an onion site. Social networking
on i2p in general could be quite damaging if malicious actors got involved. It
would be unwise, I think, to trust one GNU Social instance or another with my
http proxy destination for every eepSite I visit.

This changes the default behavior to use a different destination for every eepSite
----------------------------------------------------------------------------------

So, instead of visiting every eepSite with the same destination, we could avoid
the easiest kinds of cross-site tracking via the destination of the http proxy.
The obvious thing to do is connect to each eepSite with a unique destination,
but I considered other things. I think unique destinations is the best way to
proceed. I had considered the possibility of making requests for third-party
resources along the same destination as the eepSite requesting the third-party
resource but that seemed easier to abuse.

Attacking this New Proxy
------------------------

As I said in the top-level readme, it's possible to distinguish whether a person
is using the standard i2p http proxy or this proxy. Which means I'm a pretty
easy to spot visitor to your eepSite probably, but I don't think you can prove
actually prove it. Here's how:

### Software Fingerprinting:

The same colluding adversary can identify *this proxy* by identifying which
resources you *don't* request from their site. In other words, if a site
administrator knowingly requests off-site resources, it can identify
destinations created by this proxy by identifying destinations that *never*
request those third party resources when connecting to their eepSite. This
effect is probably worse if there are few visitors to that eepSite.

### Header Fingerprinting

As you can see from [these screenshots](SCREENSHOTS.md), the headers used to
allow 'Accept-Language' through, revealing it was in use. They don't anymore,
which is cool, but these are the only headers it filters:

        var hopHeaders = []string{
            "Proxy-Authenticate",
            "Proxy-Authorization",
            "Proxy-Connection",
            "X-Forwarded-For",
            "Accept-Language",
        }

So, if you are passing extra information, then the server will probably get that
information. Which will not only reveal that information, but also potentially
that you're using this proxy and not the default proxy if we strip out different
headers. **THIS IS WHY I STRONGLY RECOMMEND, and pre-configure if desired, a**
**Tor Browser.** Whether I think another existing project could be recommended
is a question I haven't got a good answer to at this point. Brave's new privacy
mode is a thing. Maybe browsers will get better. Maybe the web can be a little
less crappy.

### Timing Attacks on Small Anonymity Sets

If you're visiting a network of colluding eepSites attempting tracking, you're
probably currently part of a very small number of i2p users. Extremely small.
Not even my test sites are up anymore. Which means that if an eepSite requests
a remote resource from another eepSite, there's a good chance that those two
requests will be the only requests to those two eepSites contemporary to
eachother, establishing a possible link between the two destinations. To combat
this, the new proxy tears down destinations that become inactive and creates
new destinations if they become active again, but a more active attacker could
request the resource over and over again to keep the tunnel alive. This doesn't
seem to gain them much, except to establish a link between contemporary
activities on two eepSites, but it does do that. It might be possible to use
some kind of caching to prevent repeated resource requests within the lifespan
of the http proxy if it's possible to do that without causing cache-timing
attacks.

So those are the two I know about. If you know of other's, I'm pretty happy to
deal with them publicly. I've visibly insisted this is experimental for quite a
long time, hopefully nobody is misusing it. Also I'm reasonably confident it
works.

Other avenues for Destination-Isolation:
----------------------------------------

This elaborate system is designed to accomplish the goal of presenting an http
proxy server. Other approaches are probably better for things that are not web
browsers. A socksifer using the SAM bridge to generate ephemeral destinations
is probably a good idea.

Thoughts on long-term destinations for eepSites
-----------------------------------------------

