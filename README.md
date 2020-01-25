seoladh
=======

A client/server to register and serve IP addreses by name.

It consists of a very simple key-value store backed and a HTTP server. The
server listens for PUT requests at an endpoint, for example

    curl -X PUT https://seoladh.myserver.com:12345/home.laptop.com

In this case, the service is running (HTTPS) on port 11923 of a machine which
is registered as seoladh.myserver.com. The response of the call is just a 201
(Created), to indicate that the mapping has been created. Suppose the call was
made from a machine at 192.2.3.4 (which may be behind a NAT gateway). Then a
subsequent call from any machine to

    curl https://seoladh.myserver.com:12345/home.laptop.com

will result in a 200 with a Content-Type of text/plain, containing the value
"192.2.3.4". It's designed to be quick, so a machine could updated the server
every few seconds, thereby more or less guaranteeing the veracity of the
mappings.
