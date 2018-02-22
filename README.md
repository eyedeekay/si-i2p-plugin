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

### The http proxy

Still *extremely experimental*, but currently it is possible to set your web
browser's HTTP proxy to localhost:4443 and use it to browse eepSites. I just got
it working and it's not been tested much yet, YMMV.

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
pipe. To read out the most recent response, cat out the parent/recv pipe. Lastly
to close all the pipes and clean up, echo "y" into parent/del.

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

Making sure it the aggregator pipe does things correctly. I'm thinking I'll use
this to teach myself more about AppArmor and SELinux profiles next.

What the final version should do:
---------------------------------

The final version should use the parent pipe and the aggregating pipe to send
and recieve requests as an http proxy in the familiar way.

Version Roadmap:

  * ~~0.17 - Named Pipes work for top-level i2p domains and can retrieve~~
   ~~directories under a site~~
  * ~~0.18 - Named Pipes for i2p domains and can retrieve subdirectories,~~
   ~~which it caches in clearly-named folders as normal files(Containing HTML)~~
  * 0.19 - Expose an http proxy that hooks up to the existing infrastructure
   for destination isolation
  * 0.20 - Ready for more mainstream testing, should successfully isolate
   requests for resources embedded in the retrieved web pages.
  * 0.21 - First worthwhile release for people who aren't shell enthusiasts.

Silly Questions I'm asking myself about how I want it to work:
--------------------------------------------------------------

Should it do filtering? I really don't think so but if there's a simple way to
strip outgoing information then maybe. I dislike complexity. It's why this has
taken so long.

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

Installation and Usage:
=======================

This is still very much a work in progress. It's being developed on Debian, but
it can also be easily built on any system with Go available and can also easily
be part of a Docker image.

Build Dependencies: Go, go tools, a compiler and libc, gosam, git(for grabbing
gosam library), make

Optional: Docker, checkinstall(for building deb packages, will be replaced with
proper packaging when .25 is reached.

Runtime Dependencies: An i2p router with the SAM bridge enabled.

Building is accomplished with a simple

        make

*BUT please note that this is going to "go get" my copy of the gosam repo and*
*place it in your GOPATH.* But no other commands will be necessary. The
executable will be placed in working_directory/bin and can be run as a the user
on any system with the SAM bridge enabled. But if you run

        sudo make install

after, it will install into /usr/local/bin/ by default and more importantly, it
can be started as a system service running as it's own user(For now, initscripts
exist and untested systemd units exist).

On Debian or Ubuntu
-------------------

A recommended install procedure is

        sudo apt-get install git golang make
        git clone https://github.com/eyedeekay/si-i2p-plugin
        make checkinstall
        sudo dpkg -i ../si-i2p-plugin*.deb

This will allow you to use your package manager to install and uninstall the
service and keep your system aware of the package.

On Docker
---------

I do the building of the Docker image in two stages to make sure that the image
is as tiny as possible. And I mean tiny. And unprivileged and incapable too. It
can run on Scratch, as it can be statically linked(More on that in a moment) but
I can't get it to statically link on Debian with gccgo yet. What we do is build
a statically-linked version of the plugin in one container, extract it, and
generate a second container containing only the statically linked version of the
plugin itself.

Building a static version:
--------------------------

For now, building a static version of the plugin requires using Docker to build
it in an alpine container. To do this, run,

        make static

which will automatically build the Dockerfiles/Dockerfile.static with the
default settings, compile the statically-linked variant of the program, and
copy it from the container to the host.

Building the runtime image:
---------------------------

Once you have a static binary to use, you can create the minimal docker image
you'll want to use to run the program. To build this container based on the
docker Scratch image, run

        make docker

which will copy the statically-linked variant into a base container, along with
a minimal user and a statically compiled bash shell so that the process can run
unprivileged in the /opt directory of the scratch container.

Running the docker image:
-------------------------

Finally, to actually run the docker image, you can either run it manually or
use

        make docker-run

to start the container. All in all, the whole thing weigh's in at just under 10
mb. Which is at least much lighter than it would be if I had a whole Ubuntu
container in here.

Building a .deb containing the statically-linked variant:
---------------------------------------------------------

It's also possible to use the statically-linked variant as the basis for a
package by running the command

        checkinstall-static

which you can then install with

        sudo dpkg -i ../si-i2p-plugin*-static.deb

