FROM golang

WORKDIR /go/src/

RUN go get github.com/nvisibleinc/go-ari-library
RUN go get golang.org/x/net/websocket
RUN mkdir go-ari-proxy
COPY . go-ari-proxy
RUN go install go-ari-proxy

CMD /go/bin/go-ari-proxy
