doduc
=====

A dynamic DNS system for DigitalOcean

Background
----------

The DigitalOcean API allows domain name records to be updated. This means that
it is possible for DigitalOcean DNS records to be used as a dynamic DNS system.
Existing clients seem to query third-party services like
[ifconfig.me](http://ifconfig.me/ip). Rather than relying on these services, I
wanted to use my own DigitalOcean droplet to serve the IPs.

doduc has two components. The first is a server, which runs on your DigitalOcean
droplet (or any other remote server). It runs on port 18768 and responds to all
HTTP requests with their source IP. The second is a client, which runs on your
local computer. It fetches your IP address from the server and updates your
DigitalOcean DNS for your chosen domain and subdomain with the address. Both the
server and client save output to a log file (`doduc-server.log` and
`doduc-client.log`, respectively).

The server and client are designed to run continuously, so you will probably
want to run them in the background. On POSIX operating systems, you can run

`doduc-[program] &`

In Windows PowerShell, you can run

`Start-Process doduc-[program] -ArgumentList [comma-separated list of arguments] -WindowStyle Hidden`

Building
--------

Run `go build` in `doduc-client` and `doduc-server`.

Usage
-----

### Server

Simply run `doduc-server`.

### Client

If you wish to use another server to get your external IP, the only requirement
is that it must respond to a GET HTTP request with a valid IP address (e.g.
[ifconfig.me](ifconfig.me/ip)).

To run the client, you must first
[generate an OAuth token](https://cloud.digitalocean.com/settings/tokens/new)
for the client to use. Save the generated token to a file (e.g. `token` in the
root of this repository).

Now you can run the client by specifying the DigitalOcean domain you want to
update, the subdomain you wish to point to your IP address, the URL of your IP
server, and the path to the file containing your OAuth token. For instance, if
you wanted `home.mywebsite.com` to point to your IP address, you had
`doduc-server` running on `www.mywebsite.com`, and you had your OAuth token
saved in `token` in the current directory, you would run the client on your
machine as follows:

`doduc-client -domain mywebsite.com -subdomain home -ip_server http://www.mywebsite.com:18768 -token token`

By default, the address will be updated every 5 minutes. To specify the interval
in seconds, use the `interval` flag. For instance, if you wanted updates to
occur every 10 minutes:

`doduc-client -domain mywebsite.com -subdomain home -ip_server http://www.mywebsite.com:18768 -token token -interval 600`