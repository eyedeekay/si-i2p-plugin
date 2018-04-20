Installation and Usage:
=======================

This is still very much a work in progress. It's being developed on Debian, but
it can also be easily built on any system with Go available and can also easily
be part of a Docker image.

FOSS UNIX-Like
==============

Build Dependencies: Go, go tools, a compiler and libc, gosam, git(for grabbing
gosam library), make

Optional: Docker, checkinstall(for building deb packages, will be replaced with
proper packaging when .25 is reached.

Runtime Dependencies: An i2p router with the SAM bridge enabled.

Building is accomplished with a simple

        make build

and

        make rebuild

will build a fresh version.

*BUT please note that this is going to "go get" my copy of the gosam repo and*
*place it in your GOPATH.* But no other commands will be necessary. The
executable will be placed in working_directory/bin and can be run as a the user
on any system with the SAM bridge enabled. But if you run

        sudo make install

after, it will install into /usr/local/bin/ by default and more importantly, it
can be started as a system service running as it's own user(For now, initscripts
exist and untested systemd units exist).

From this point on, all the executables are statically-linked and
position-independent.

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

Building the runtime image:
---------------------------

No more complicated shit. It's all packaged. just do:

        make docker-setup

Note that the Docker container will use 44443 instead of 4443.


OSX
===

**UNTESTED** An untested OSX static binary is in the repository at
bin/si-i2p-plugin.bin. It's a free-standing, unpackaged mach-O static
executable. But it's built from identical code to the Linux version.

To cross-compile for OSX, you can do:

        make osx

or if you're on a Mac, you can do:

        make native

Windows
=======

**UNTESTED** An untested Windows static binary is in the repository at
bin/si-i2p-plugin.exe. It's just a Windows PE executable, and I'm pretty sure
it's a static executable. It does a few things differently than the Unix-ish
versions, in particular it treats the named pipes as regular old files. In the
case of "input" pipes, they are manually truncated to 0 length afte inputs are
processed.

To cross-compile for Windows, you can do:

        make windows

or if you're on a Windows PC, you can do:

        go build -a -o bin/si-i2p-plugin.exe ./src

from the project root directory.
