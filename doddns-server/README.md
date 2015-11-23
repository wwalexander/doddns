doddns-server
=============

Building
--------

    go build

Usage
-----

    doddns-server [-cert path] [-key path] [-port port]

doddns-server responds to HTTP requests with their source IP address.
doddns-server runs on port 18768, or the port named by the -port flag. If the
-cert and -key flags are specified, doddns-server will listen for HTTPS
connections; otherwise, doddns-server will listen for HTTP connections.
