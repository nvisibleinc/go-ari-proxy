go-ari-proxy
============

An implementation of the [go-ari-library][1] used to connect to the Asterisk
REST Interface for delivery of Events and Commands across a message bus.

* publishes Events from ARI websocket to a remote client via message bus
* consumes Commands from a client using a message bus and POSTs them to Asterisk
* publishes Command responses back onto the message bus for the client

Overview
--------
The `go-ari-proxy` is an application that makes it easier to build external applications for the [Asterisk](http://github.com/asterisk) REST Interface (ARI). Primarily, we wrap the ARI messages in a structure that can be easily disseminated and sent across a message bus using `json`. From there, the application can decapsulate the message, act on the ARI portion, and send back a response in a similar manner.

How It Works
------------
The `go-ari-proxy` application works by using a message bus (NATS, NSQ, RabbitMQ) and having a _broadcast_ channel (topic, queue) that applications listen to. When a `StasisStart` is recieved by the proxy, we perform an `AppStart` which broadcasts this information. The message bus then distributes this to an application listening on the topic. From there, the application knows what topic to subscribe to in order to continue the conversation. Currently this setup channel is unidirectional, however in the future it may be further enhanced/designed to allow for applications to acknowledge setup of a dialog.

The primary purpose of the broadcast channel is to distribute to one or more application instances, and then a dialog is created over three (3) additional topics. These topics are: *Events*, *Commands*, and *CommandResponses*.

![Application Topic Distribution](docs/images/application-topic-distribution.jpg "Application Topic Distribution")

Installation
------------
```
$ cd $GOPATH/src
$ git clone https://github.com/nvisibleinc/go-ari-proxy
$ go install go-ari-proxy
```

Configuration
--------------

Docker Container
----------------



Licensing
---------
`go-ari-proxy` is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/docker/docker/blob/master/LICENSE) for the full
license text.

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

   [1]: https://github.com/nvisibleinc/go-ari-proxy
