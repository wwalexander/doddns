doddns-client
=============

Building
--------

    go build

Usage
-----

    doddns-client [-interval seconds] [domain] [subdomain] [server] [token]

doddns-client periodically updates the DNS A record for subdomain.domain using
the IP address returned by the named server via HTTP. To authenticate,
doddns-client uses the DigitalOcean API token saved at the named path. By
default, doddns-client updates the record every 5 minutes; this interval can be
set to a given number of second using the -interval flag.

doddns-server provides an implementation of a server that can be used by
doddns-client.
