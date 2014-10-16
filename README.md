go-ari-proxy
============

An implementation of the `go-ari-library`[1] used to connect to the Asterisk
REST Interface for delivery of Events and Commands across a message bus.

* publishes Events from ARI websocket to a remote client via message bus
* consumes Commands from a client using a message bus and POSTs them to Asterisk
* publishes Command responses back onto the message bus for the client

Installation
------------
```go
$ go build go-ari-proxy
```

Licensing
---------
> Copyright 2014 N-Visible Technology Lab, Inc.
> 
> This program is free software; you can redistribute it and/or
> modify it under the terms of the GNU General Public License
> as published by the Free Software Foundation; either version 2
> of the License, or (at your option) any later version.
> 
> This program is distributed in the hope that it will be useful,
> but WITHOUT ANY WARRANTY; without even the implied warranty of
> MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
> GNU General Public License for more details.
> 
> You should have received a copy of the GNU General Public License
> along with this program; if not, write to the Free Software
> Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.

[1] https://github.com/nvisibleinc/go-ari-proxy