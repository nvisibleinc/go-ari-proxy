go-ari-library
==============
A library for building an Asterisk REST Interface proxy and client using a
message bus backend for delivery of messages, written in the Go programming
language.

This library abstracts the message bus from the application by providing an
interface for setting up channels to consume Events and Commands in a bus
agnostic way.

Installation
------------
```
$ go import https://github.com/nvisibleinc/go-ari-library
```

Usage
-----
```go
import (
	"https://github.com/nvisibleinc/go-ari-library"
)
```

For a useful example of usage of this library, see the [go-ari-proxy][1] and
[ari-voicemail][2] projects.

Licensing
---------
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

   [1]: https://github.com/nvisibleinc/go-ari-proxy
   [2]: https://github.com/nvisibleinc/ari-voicemail
