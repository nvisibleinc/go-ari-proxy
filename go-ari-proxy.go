package main

import (
	"encoding/json"
	"net/http"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"bytes"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"time"

	"code.google.com/p/go.net/websocket"
	"go-ari-library"
)

// Var config contains a Config struct to hold the proxy configuration file.
var (
	config Config
	client = &http.Client{}
	proxyInstances map[string] *proxyInstance
)

// Struct Config holds the configuration for the proxy.
// The Config struct contains the information the was unmarshaled from the
// configuration file for ths proxy.
type Config struct {
	Origin        string      `json:"origin"`		// connection to ARI events
	ServerID      string      `json:"server_id"`	// unique server ident
	Applications  []string    `json:"applications"`	// slice of applications to listen for
	Websocket_URL string      `json:"websocket_url"`// websocket to connect to
	Stasis_URL    string      `json:"stasis_url"`   // Base URL of ARI REST API
	WS_User       string      `json:"ws_user"`		// username of websocket connection
	WS_Password   string      `json:"ws_password"`	// pass of websocket connection
	MessageBus    string      `json:"message_bus"`	// type of message bus to publish to
	BusConfig     interface{} `json:"bus_config"`	// configuration of the message bus we're publishing to
}

// ProxyInstance struct contains the channels necessary for communications
// to/from the various message bus topics and the event channel. This is
// primarily used as the communications bus for setting up new instances of
// applications.
type proxyInstance struct {
	commandChannel		chan []byte
	responseChannel		chan []byte
	Events				chan []byte
	quit				chan int
	ariObjects			[]string
}

// eventInfo struct contains the information about an event that comes in.
// Information about the event that we need to make a determination on the proxy side.
// Track information associated with a given application instance.
type eventInfo struct {
	Type      string	`json:"type"`		// event type
	Application string	`json:"application"`
	Bridge    minBridge `json:"bridge"`		// bridge ID
	Channel   minChan   `json:"channel"`	// channel ID
	Recording minRec	`json:"recording"`	// recording name	TODO: why no id?
	Playback  minPlay	`json:"playback"`	// playback ID
}

// minPlay struct is a minimal struct that contains the Playback ID.
type minPlay struct {
	ID	string	`json:"id"`
}

// minRec struct is a minimal struct that contains the Recording Name (no ID).
type minRec struct {
	Name	string	`json:"name"`
}

// minBridge struct is used to get the ID of a bridge in an event.
type minBridge struct {
	ID	string	`json:"id"`
}

// minChan struct is used to get the ID of a channel in an event.
type minChan struct {
	ID	string	`json:"id"`
}

// ID struct is used to get the ID or name of objects returned from ARI REST calls.
type ID struct {
	ID		string `json:"id"`
	Name	string `json:"name"`
}

// Init parses the configuration file by unmarshaling it into a Config struct.
func init() {
	var err error

	// parse the configuration file and get data from it
	configpath := flag.String("config", "./config.json", "Path to config file")

	flag.Parse()
	configfile, err := ioutil.ReadFile(*configpath)
	if err != nil {
		log.Fatal(err)
	}
	// read in the configuration file and unmarshal the json, storing it in 'config'
	json.Unmarshal(configfile, &config)
	fmt.Println(&config)
	proxyInstances = make(map[string] *proxyInstance)
}

// PublishMessage takes an ARI event from the websocket and places it on the
// producer channel.
// Accepts two arguments:
// * a string containing the ARI message
// * a producer channel
func PublishMessage(ariMessage string, producer chan []byte) {
	// unmarshal into an ari.Event so we can append some extra information
	var info eventInfo
	var message ari.Event
	var pi *proxyInstance
	var exists bool = false
	json.Unmarshal([]byte(ariMessage), &message)
	json.Unmarshal([]byte(ariMessage), &info)
	message.ServerID = config.ServerID
	message.Timestamp = time.Now()
	message.ARI_Body = ariMessage
	switch  {
	case info.Type == "StasisStart":
		// since we're starting a new application instance, create the proxy side
		dialogID := ari.UUID()
		fmt.Println("Dialog ID is:", dialogID)
		as, err := json.Marshal(ari.AppStart{Application: info.Application, DialogID: dialogID})
		producer <- as
		time.Sleep(50 * time.Millisecond)
		if err != nil {
			return
		}
		pi = initProxyInstance(dialogID)
		proxyInstances[info.Channel.ID] = pi
		exists = true
	case info.Type == "StasisEnd":
		// on application end, perform clean up checks
		pi, exists = proxyInstances[info.Channel.ID]
		if exists {
			//pi.removeAllObjects()
		}
	case info.Type == "BridgeDestroyed":
		pi, exists = proxyInstances[info.Bridge.ID]
		if exists {
			pi.removeObject(info.Bridge.ID)
		}
	case info.Type == "ChannelDestroyed":
		pi, exists = proxyInstances[info.Channel.ID]
		if exists {
			pi.removeObject(info.Channel.ID)
		}
	case strings.HasPrefix(info.Type, "Channel"):
		pi, exists = proxyInstances[info.Channel.ID]
	case strings.HasPrefix(message.Type, "Bridge"):
		pi, exists = proxyInstances[info.Bridge.ID]
	default:
		fmt.Println("No handler for event type")
		//pi, exists = proxyInstances[]
		// if not matching, then we need to perform checks against the
		// existing map to determine where to send this ARI message.
	}
	// marshal the message back into a string
	busMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("[DEBUG] Bus Data:\n%s\n", busMessage)

	// push the busMessage onto the producer channel
	if exists {
		pi.Events <- busMessage
	}
}

// initProxyInstance initializes a new proxy instance.
func initProxyInstance(dialogID string) *proxyInstance {
	var p proxyInstance
	p.quit = make(chan int)
	p.Events = ari.InitProducer(strings.Join([]string{"events", dialogID}, "_"))
	go p.runCommandConsumer(dialogID)
	return &p
}

// runEventHandler sets up a websocket connection to an ARI application.
func runEventHandler(s string, producer chan []byte) {
	// Connect to the websocket backend (ARI)
	var ariMessage string
	url := strings.Join([]string{config.Websocket_URL, "?app=", s, "&api_key=", config.WS_User, ":", config.WS_Password}, "")
	ws, err := websocket.Dial(url, "ari", config.Origin)
	if err != nil {
		log.Fatal(err)
	}

	// Start the producer loop. Every message received from the websocket is
	// passed to the PublishMessage() function.
	for {
		err = websocket.Message.Receive(ws, &ariMessage)		// accept the message from the websocket
		if err != nil {
			log.Fatal(err)
		}
		go PublishMessage(ariMessage, producer)	// publish message to the producer channel
	}
}

// shutDown closes the quit channel to signal all of a ProxyInstance's goroutines
// to return
func (p *proxyInstance) shutDown () {
	close(p.quit)
}

func (p *proxyInstance) addObject(id string) {
	for i:= range p.ariObjects {
		if p.ariObjects [i] == id {
			//object already is associated with this proxyInstance
			return
		}
	}
	p.ariObjects = append(p.ariObjects, id)
	proxyInstances[id] = p
	
}
func (p *proxyInstance) removeObject(id string) {
	// remove an object from the map.
	for i := range p.ariObjects {
		if p.ariObjects[i] == id {
			// rewrite the p.ariObjects string slice to append all values up to
			// the index value of 'i', and all values of 'i'+1 and later.
			p.ariObjects = append(p.ariObjects[:i], p.ariObjects[i+1:]...)
		}
	}

	// remove the instance from our tracking map
	delete(proxyInstances, id)

	// if there are no more objects, shut'rdown
	if len(p.ariObjects) == 0 {
		p.shutDown()
	}
}

func (p *proxyInstance) removeAllObjects() {
	// remove all objects from the map as our application is shutting down.
	for _, obj := range p.ariObjects {
		delete(proxyInstances, obj)
	}
	p.shutDown()	// destroy the application / proxy instance
}

// runCommandConsumer starts the consumer for accepting Commands from
// applications.
func (p *proxyInstance) runCommandConsumer(dialogID string) {
	commandTopic := strings.Join([]string{"commands", dialogID}, "_")
	responseTopic := strings.Join([]string{"responses", dialogID}, "_")
	fmt.Println("Topics are:", commandTopic, " ", responseTopic)
	p.responseChannel = ari.InitProducer(responseTopic)
	select {
	case <- ari.TopicExists(commandTopic):
		p.commandChannel = ari.InitConsumer(commandTopic)
	case <- time.After(10 * time.Second):
		p.removeAllObjects()
		return
	}
	
	for {
		select {
		case jsonCommand := <- p.commandChannel:
			go p.processCommand(jsonCommand, p.responseChannel)
		case <- p.quit:
			return
		}
	}
}

// processCommand processes commands from applications and submits them to the
// REST interface.
func (p *proxyInstance) processCommand(jsonCommand []byte, responseProducer chan []byte) {
	var c ari.Command
	var r ari.CommandResponse
	i := ID{ID:"", Name:""}
	fmt.Printf("jsonCommand is %s\n", string(jsonCommand))
	json.Unmarshal(jsonCommand, &c)
	fullURL := strings.Join([]string{config.Stasis_URL, c.URL, "?api_key=", config.WS_User, ":", config.WS_Password }, "")
	fmt.Println(fullURL)
	req, err := http.NewRequest(c.Method, fullURL, bytes.NewBufferString(c.Body))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	fmt.Printf("Response body is %s\n", buf.String())
	json.Unmarshal(buf.Bytes(), &i)
	if i.ID != "" {
		p.addObject(i.ID)
	} else if i.Name != "" {
		p.addObject(i.Name)
	}
	r.ResponseBody = buf.String()
	r.StatusCode = res.StatusCode
	sendJSON, err := json.Marshal(r)
	if err !=nil {
		fmt.Println(err)
	}
	fmt.Printf("sendJSON is %s\n", string(sendJSON))
	responseProducer <- sendJSON
}

// signalCatcher is a function to allows us to stop the application through an
// operating system signal.
func signalCatcher() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	sig := <-ch
	log.Printf("Signal received: %v", sig)
	os.Exit(0)
}

func main() {
	// Setup a new Event producer and Command consumer for every application
	// we've configured in the configuration file.
	ari.InitBus(config.MessageBus, config.BusConfig)
	for _, app := range config.Applications {
		producer := ari.InitProducer(app)	// Initialize a new producer channel using the ari.IntProducer function.
		go runEventHandler(app, producer)	// create new websocket connection for every application and pass the producer channel
	}

	go signalCatcher()	// listen for os signal to stop the application
	select {}
}
