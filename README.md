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

Building
--------

Run `go build` in `doduc-client` and `doduc-server`.

Usage
-----

Both the server and client save output to a log file (`doduc-server.log` and `doduc-client.log`, respectively).

### Server

doduc has a server, which runs on your DigitalOcean droplet (or any other remote
server). The server gets the source IP from incoming HTTP requests and returns
it to the client. It runs on port 18768. Thus, to run the server on  on your
machine, you would run:

`doduc-server`

If you wanted to run the server in the background so that you could exit the
shell, you could run:

`doduc-server &`

### Client

If you wish to use another server to get your external IP, the only requirement
is that it must respond to a GET HTTP request with a valid IP address.

To run the client, you must first
[generate an OAuth token](https://cloud.digitalocean.com/settings/tokens/new)
for the client to use. Save the generated token to a file (e.g. `token` in the
root of this repository).

Now you can run the client by specifying the DigitalOcean domain you want to
update, the subdomain you wish to point to your IP address, the URL of your IP
server, and the path to the file containing your OAuth token. For instance, if I
wanted `home.mywebsite.com` to point to my IP address, I had `doduc-server`
running on `www.mywebsite.com`, and I had my OAuth token saved in `token` in the
current directory, I would run the client on my machine as follows:

`doduc-client -domain mywebsite.com -subdomain home -ip_server http://www.mywebsite.com:18768 -token token`

By default, the address will be updated every 5 minutes. To specify the interval
in seconds, use the `interval` flag. For instance, if I wanted updates to occur
every 10 minutes:

`doduc-client -domain mywebsite.com -subdomain home -ip_server http://www.mywebsite.com:18768 -token token -interval 600`

Similarly to the server, the client can be run in the background. If you are running the client on Windows, you can use the `Start-Process` cmdlet:

`Start-Process doduc-client.exe -ArgumentList -domain, mywebsite.com, -subdomain, home, -ip-server, http://www.mywebsite.com:18768 -token token -WindowStyle Hidden`
