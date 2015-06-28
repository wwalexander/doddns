doduc
=====

A dynamic DNS system for DigitalOcean

Background
----------

The DigitalOcean API allows domain name records to be updated. This means that it is possible for DigitalOcean DNS records to be used as a dynamic DNS system. Existing clients seem to query third-party services like [ifconfig.me](ifconfig.me/ip). Rather than relying on these services, I wanted to use my own DigitalOcean droplet to serve the IPs.

Usage
-----

doduc has a server, which runs on your DigitalOcean droplet (or any other remote server). The server gets the source IP from the incoming HTTP requests and returns it to the client. It runs on port 26992. Run the `server` executable on the server.
