doddns-client
=============

Building
--------

    go build

Usage
-----

    doddns-client [domain] [subdomain] [server] [token]

doddns-client periodically updates the DNS A record for subdomain.domain using
the IP address returned by the named server via HTTP. To authenticate,
doddns-client uses the DigitalOcean API token saved at the named path.

doddns-server provides an implementation of a server that can be used by
doddns-client.
