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
                            send (Named Pipe)
                                 <- echo "desired-url" > destination_url.i2p/send
                            recv (Output File)
                                 <- cat destination_url.i2p/recv
                            name (Output file)
                                 <- cat destination_url.i2p/name
                            del (Named pipe)
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
