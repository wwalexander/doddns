doddns
======

A dynamic DNS system for DigitalOcean

Background
----------

The DigitalOcean API allows domain name records to be updated. This means that
it is possible for DigitalOcean DNS records to be used as a dynamic DNS system.
Existing clients seem to query third-party services like
[ifconfig.me](http://ifconfig.me/ip). Rather than relying on these services, I
wanted to use my own DigitalOcean droplet to serve the IPs.

doddns has two components. The first is a server, which runs on your DigitalOcean
droplet (or any other remote server). It runs on port 18768 by default and
responds to all HTTP requests with their source IP. The second is a client,
which runs on your local computer. It fetches your IP address from the server
and updates your DigitalOcean DNS for your chosen domain and subdomain with the
address. Both the server and client save output to a log file.

The server and client are designed to run continuously, so you will probably
want to run them in the background. On POSIX operating systems, you can run

    doddns-[program] &

In Windows PowerShell, you can run

    Start-Process doddns-[program] -ArgumentList [comma-separated list of arguments] -WindowStyle Hidden

Server
------

### Building

    go build

### Usage

    doddns-server [OPTIONS]

#### Flags

`-port`: the port to listen on (defaults to 18768)

Client
------

### Building

    go build

### Usage

Periodically update the `A` record for `SUBDOMAIN`.`DOMAIN` using the IP
returned from `SERVER`, using the DigitalOcean API token stored in `TOKEN`:

    doddns-client [OPTIONS] [DOMAIN] [SUBDOMAIN] [SERVER] [TOKEN]

If you wish to use another server to get your external IP, the only requirement
is that it must respond to a GET HTTP request with a valid IP address (e.g.
[ifconfig.me](http://ifconfig.me/ip)).

To run the client, you must first
[generate an OAuth token](https://cloud.digitalocean.com/settings/tokens/new)
for the client to use. Save the generated token to a file (e.g. `token` in the
root of this repository).

#### Options

`-interval`: the interval between updates in seconds (defaults to 300)
