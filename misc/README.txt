This is a very simple in-memory open tracker, wrapped into an I2P plugin.

The plugin starts a new http server tunnel, eepsite, and Jetty server running at port 7662.
The tracker status is available at http://127.0.0.1:7662/tracker/ .
If other files are desired on the eepsite, they can be added at eepsite/docroot .

The open tracker code and jsps were written from scratch, but depend on some code
in i2psnark.jar from the I2P installation for bencoding, and of course
on other i2p libraries.
See the license files in I2P for i2p and i2psnark licenses.
There is also some code modified from Jetty 5.1.15.
See LICENSES.txt for the zzzot and Jetty licenses.

I2P source must be installed and built in ../i2p.i2p to compile this package.

Sure, as a standalone program in its own JVM with Jetty, this would be a pig -
you should use the C opentracker instead. But since you're already running
the JVM and Jetty, running this in the same JVM probably doesn't hog to much more memory.

Valid announce URLs:
	/a
	/announce
	/announce.jsp
	/announce.php
	/tracker/a
	/tracker/announce
	/tracker/announce.jsp
	/tracker/announce.php

Valid scrape URLs:
	/scrape
	/scrape.jsp
	/scrape.php
	/tracker/scrape
	/tracker/scrape.jsp
	/tracker/scrape.php

The tracker also responds to seedless queries at
	/Seedless/index.jsp

You may use the rest of the eepsite for other purposes, for example you
may place torrent files in eepsite/docroot/torrents.
