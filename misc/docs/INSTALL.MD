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
