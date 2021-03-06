doddns
======

A dynamic DNS system for DigitalOcean

Background
----------

The DigitalOcean API allows domain name records to be updated. This means that
it is possible for DigitalOcean DNS records to be used as a dynamic DNS system.
Existing clients seem to query third-party services like
http://ifconfig.me/ip. Rather than relying on these services, I
wanted to use my own DigitalOcean droplet to serve the IPs.

doddns is composed of a client and server, which run on the local machine and a
remote server respectively. The server is meant to run on a DigitalOcean droplet.
It responds to HTTP requests with the IP address of the client a la
http://ifconfig.me/ip. The client updates the DNS record for the given domain
and subdomain with this address automatically.