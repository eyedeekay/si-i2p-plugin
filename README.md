Destination-Isolating i2p HTTP Proxy(SAM Application)
=====================================================

This is an i2p SAM application which presents an HTTP proxy(on port 4443 by
default) that acts as an intermediate between your browser and the i2p network.
Then it uses the SAM library to create a unique destination for each i2p site
that you visit. This way, your unique destination couldn't be used to track you
with a network of colluding sites. I doubt it's a substantial problem right now
but it might be someday. Facebook has an onion site, and i2p should have
destination isolation before there is a facebook.i2p.

Excitingly, after like a year of not being able to devote the time I should
have to this, I've finally made it work. I am successfully using this proxy to
browse eepSites pretty normally. There are still a few bugs to chase down, but
if somebody else wanted to try it out, I might not *totally* embarass myself.

What works so far:
------------------

### The http proxy

Still *extremely experimental*, but currently it is possible to set your web
browser's HTTP proxy to localhost:4443 and use it to browse eepSites. I just got
it working and it's not been tested much yet, YMMV.

I think I've solved the race condition thing now, which makes it usable. It's
not even that bad. As expected, it takes a little longer to get the site in the
first place.

#### Examples

##### curl

        curl -x 127.0.0.1:4443 http://i2p-projekt.i2p

##### surf

        export http_proxy="http://127.0.0.1:4443"
        surf http://i2p-projekt.i2p

#### Current Concerns:

URL validation needs quite a bit of work, and I may need some additional way of
validating responses.

I haven't been able to observe any DNS leaks yet, but that doesn't mean they
aren't there. My plan is to implement some kind of proper URL validation for it.

Before version 0.21, a framework for generating service tunnels ad-hoc will also
be in place. This will be used for fuzz-testing the http proxy and the pipe
proxy. Almost everything will be improved by the availability of this.

### The pipes

It currently functions well as a file/pipe based interface to http services on
the i2p network. It doesn't work as an http proxy yet.

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
pipe(one is loaded automatically right now, but will be removed in the future).
To read out the most recent response, cat out the parent/recv pipe. Lastly, to
close all the pipes and clean up, echo "y" into parent/del.

Behind that, there is a system which uses named pipes to allow a user to send
and recieve requests and get information about eepSites on the i2p network. If
you were to, for instance, make a request for i2p-projekt.i2p through
parent/send, it would look for the SAM session associated with that site(or
create one if it doesn't exist) in a folder called "i2p-projekt.i2p". Inside
that folder will be 5 files corresponding to the named pipes and the output
files:

        destination_url.i2p/
                            send
                                 <- echo "desired-url" > destination_url.i2p/send
                            recv (Output File)
                                 <- cat destination_url.i2p/recv
                            name (Named pipe but will probably become an output file)
                                 <- cat destination_url.i2p/name
                            del
                                 <- echo "y" > destination_url.i2p/del
                            time (Output File)
                                 <- cat destination_url.i2p/time

In order to use them to interact with eepSites, you may either make your
requests to the parent pipes  which will delegate the responses to the child
pipes automatically, or you may manually pipe the destination URL into
destination\_url.i2p/send, and pipe out the result from
destination\_url.i2p/recv. To retrieve the full cryptographic identifier of the
eepSite, pipe out the destination from destination\_url.i2p/name and to close
the pipe, pipe anything at all into destination\_url.i2p/del. The final field,
destination\_url.i2p/time is the time which the page in the folder was last
recieved.

When you retrieve a sub-directory of a site or a URL under the domain, a new set
of named pipes and output files will be created in a directory corresponding
to that URL underneath the destination\_url.i2p/ folder. These folders can
be created using either the parent/send pipe, which will automatically route
it through the correct destination, or through destination\_url.i2p/send which
will send it through a specific destination. The final behavior of this pipe is
not yet determined but may be modified to only allow requests to the already
authorized destination or not, as a way of electively sharing information
between eepSites if so desired. For now, no validation of the intended
destination is done in the child proxies. A subdirectory managed by a child
proxy will look like

        destination_url.i2p/
                            subdirectory/
                                         recv
                                            cat destination_url.i2p/subdirecctory/recv
                                         time
                                            cat destination_url.i2p/subdirectory/time
                                         del
                                            echo "y" > destination_url.i2p/subdirectory/del

Note that the send ane name pipes are not present as they are provided by the
managing child proxy.

Also, caching, after a fashion, is already available because the recieved files
are just files.

What I'm doing right now:
-------------------------

Improving i2p url pre-validation and implementing pipe-controlled service
tunnels.

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
  * 0.20 - Ready for more mainstream testing, should successfully isolate
   requests for resources embedded in the retrieved web pages and should be able
   to generate services on the fly by talking to the SAM bridge.
  * 0.21 - First worthwhile release for people who aren't shell enthusiasts.

Silly Questions I'm asking myself about how I want it to work:
--------------------------------------------------------------

Should it do filtering? I really don't think so but if there's a simple way to
strip outgoing information then maybe. I dislike complexity. It's why this has
taken so long.

What doesn't it do?
===================

Much filtering. It sanitizes some unnecessary headers, but doesn't filter
javascript. It will never filter Javascript on it's own, but it will make some
attempt to filter headers and it now rewrite's the user-agent string(with the
default being the string offered by the default i2p http proxy.

Installation and Usage:
=======================

Moved to [misc/docs/INSTALL.md](https://github.com/eyedeekay/si-i2p-plugin/tree/master/misc/docs/INSTALL.md)
