go-ari-proxy
============
An implementation of the [go-ari-library][1] used to connect to the Asterisk
REST Interface for delivery of Events and Commands across a message bus.

* publishes Events from ARI websocket to a remote client via message bus
* consumes Commands from a client using a message bus and POSTs them to Asterisk
* publishes Command responses back onto the message bus for the client

## Overview
The `go-ari-proxy` is an application that makes it easier to build external applications for the [Asterisk](http://github.com/asterisk) REST Interface (ARI). Primarily, we wrap the ARI messages in a structure that can be easily disseminated and sent across a message bus using `json`. From there, the application can decapsulate the message, act on the ARI portion, and send back a response in a similar manner.

## How It Works

The `go-ari-proxy` application works by using a message bus (NATS, RabbitMQ) and having a _broadcast_ channel (topic, queue) that applications listen to. When a `StasisStart` is recieved by the proxy, we perform an `AppStart` which broadcasts this information.

The message bus then distributes this to an application listening on the topic. From there, the application knows what topic to subscribe to in order to continue the conversation. Currently this setup channel is unidirectional, however in the future it may be further enhanced/designed to allow for applications to acknowledge setup of a dialog.

The primary purpose of the broadcast channel is to distribute to one or more application instances, and then a dialog is created over three (3) additional topics. These topics are: *Events*, *Commands*, and *CommandResponses*.

NOTE:  Prior revisions of the go-ari-library upon which this proxy is built had support for the NSQ message bus.  Unfortunately, as a result of the eventually consistent semantics of NSQ's topic discovery mechanism, its use is no longer supported and the code was removed.

### Application topic
![Application Topic](docs/images/app-topic.jpg "Application Topic")

1. proxy starts up and connects to the websocket
2. proxy connects to the signalling topic, named for the application (configured via JSON)
3. applications connect to same signalling topic based on their configuration
4. the message bus distributes messages in a round-robin fashion to application (message bus dependent)

### Dialog topic
![Dialog Topic](docs/images/dialog-topics.jpg "Application Topic")

1. on `StasisStart` event, the proxy creates a new _ProxyInstance_, creating a unique _dialogID_, and connect to three new topics
2. On new _dialog_ setup, the proxy sends a new `AppStart` event across the signalling channel to tell the application which topics to listen for _Events_, to send _Commands_, and to listen for _Command Responses_.

### Application topic distribution
![Application Topic Distribution](docs/images/application-topic-distribution.jpg "Application Topic Distribution")

## Installation
```
$ cd $GOPATH/src
$ go get https://github.com/nvisibleinc/go-ari-library
$ git clone https://github.com/nvisibleinc/go-ari-proxy
$ go install go-ari-proxy
```

## Configuration
TODO

## Docker Container
TODO

## Licensing
> Copyright 2015 N-Visible Technology Lab, Inc.
> 
> Licensed under the Apache License, Version 2.0 (the “License”); you may not
> use this file except in compliance with the License. You may obtain a copy
> of the License at
> 
> http://www.apache.org/licenses/LICENSE-2.0
> 
> Unless required by applicable law or agreed to in writing, software distributed
> under the License is distributed on an “AS IS” BASIS, WITHOUT WARRANTIES OR
> CONDITIONS OF ANY KIND, either express or implied. See the License for the
> specific language governing permissions and limitations under the License.

