package main

import (
	"github.com/nvisibleinc/go-ari-library"
	"strings"
	"sync"
)

// proxyInstanceMap is a singleton which holds the map
// of active proxy instances.
type proxyInstanceMap struct {
	instanceMap map[string]*proxyInstance
	mapLock     *sync.RWMutex
}

// NewproxyInstanceMap initializes a new proxy mapping.
func NewproxyInstanceMap() *proxyInstanceMap {
	p := proxyInstanceMap{}
	p.mapLock = new(sync.RWMutex)
	p.instanceMap = make(map[string]*proxyInstance)
	return &p
}

// Config holds the configuration for the proxy.
// The Config struct contains the information the was unmarshaled from the
// configuration file for ths proxy.
type Config struct {
	Origin        string      `json:"origin"`        // connection to ARI events
	ServerID      string      `json:"server_id"`     // unique server ident
	Applications  []string    `json:"applications"`  // slice of applications to listen for
	Websocket_URL string      `json:"websocket_url"` // websocket to connect to
	Stasis_URL    string      `json:"stasis_url"`    // Base URL of ARI REST API
	WS_User       string      `json:"ws_user"`       // username of websocket connection
	WS_Password   string      `json:"ws_password"`   // pass of websocket connection
	MessageBus    string      `json:"message_bus"`   // type of message bus to publish to
	BusConfig     interface{} `json:"bus_config"`    // configuration of the message bus we're publishing to
}

// proxyInstance struct contains the channels necessary for communications
// to/from the various message bus topics and the event channel. This is
// primarily used as the communications bus for setting up new instances of
// applications.
type proxyInstance struct {
	commandChannel  chan []byte
	responseChannel chan []byte
	Events          chan []byte
	quit            chan int
	ariObjects      []string
}

// NewProxyInstance initializes a new proxy instance.
func NewProxyInstance(dialogID string) *proxyInstance {
	var p proxyInstance
	p.quit = make(chan int)
	p.Events = ari.InitProducer(strings.Join([]string{"events", dialogID}, "_"))
	go p.runCommandConsumer(dialogID)
	return &p
}

// eventInfo struct contains the information about an event that comes in.
// Information about the event that we need to make a determination on the proxy side.
// Track information associated with a given application instance.
type eventInfo struct {
	Type        string    `json:"type"` // event type
	Application string    `json:"application"`
	Bridge      minBridge `json:"bridge"`    // bridge ID
	Channel     minChan   `json:"channel"`   // channel ID
	Recording   minRec    `json:"recording"` // recording name	TODO: why no id?
	Playback    minPlay   `json:"playback"`  // playback ID
}

// minPlay struct is a minimal struct that contains the Playback ID.
type minPlay struct {
	ID string `json:"id"`
}

// minRec struct is a minimal struct that contains the Recording Name (no ID).
type minRec struct {
	Name string `json:"name"`
}

// minBridge struct is used to get the ID of a bridge in an event.
type minBridge struct {
	ID string `json:"id"`
}

// minChan struct is used to get the ID of a channel in an event.
type minChan struct {
	ID string `json:"id"`
}

// ID struct is used to get the ID or name of objects returned from ARI REST calls.
type ID struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
