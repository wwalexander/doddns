doddns-client
=============

Building
--------

    go build

Usage
-----

    doddns-client [-type type] domain subdomain URI token

doddns-client periodically updates the DNS record for subdomain.domain using
the IP address returned by the named HTTP/HTTPS URI. To authenticate,
doddns-client uses the DigitalOcean API token saved at the named path. The
record is updated to match the TTL of the domain, in order to avoid doing
useless DNS updates between TTL timeouts. By default, doddns-client will update
the A record; arbitrary record types can be specified with the -type flag.