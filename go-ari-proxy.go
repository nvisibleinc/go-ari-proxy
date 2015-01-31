package main

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"flag"
	"github.com/nvisibleinc/go-ari-library"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Var config contains a Config struct to hold the proxy configuration file.
var (
	config         Config            // main proxy configuration structure
	client         = &http.Client{}  // connection for Commands to ARI
	proxyInstances *proxyInstanceMap // maps the per-dialog proxy instances
	Debug          *log.Logger
	Info           *log.Logger
	Warning        *log.Logger
	Error          *log.Logger
)

// signalCatcher is a function to allows us to stop the application through an
// operating system signal.
func signalCatcher() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	sig := <-ch
	log.Printf("Signal received: %v", sig)
	os.Exit(0)
}

// Init parses the configuration file by unmarshaling it into a Config struct.
func init() {
	var err error

	// Setup our logging interfaces
	Debug = ari.InitLogger(os.Stdout, "DEBUG")
	Info = ari.InitLogger(os.Stdout, "INFO")
	Warning = ari.InitLogger(os.Stdout, "WARNING")
	Error = ari.InitLogger(os.Stderr, "ERROR")

	// parse the configuration file and get data from it
	Info.Println("Loading configuration for proxy.")
	configpath := flag.String("config", "./config.json", "Path to config file")

	Debug.Println("Parsing the configuration file.")
	flag.Parse()
	Debug.Println("Reading in the configuration from the file.")
	configfile, err := ioutil.ReadFile(*configpath)
	if err != nil {
		Error.Fatal(err)
	}
	// read in the configuration file and unmarshal the json, storing it in 'config'
	Debug.Println("Unmarshaling proxy configuration.")
	json.Unmarshal(configfile, &config)
	Debug.Println(&config)
	Debug.Println("Initialize the proxy instance map.")
	proxyInstances = NewproxyInstanceMap() // initialize a new proxy instance map
}

func main() {
	// Setup a new Event producer and Command consumer for every application
	// we've configured in the configuration file.
	Info.Println("Initializing the message bus.")
	ari.InitBus(config.MessageBus, config.BusConfig)
	for _, app := range config.Applications {
		/*
			Create a new producer which is responsible for the initial topic on the message bus which is used
			to signal the setup of per-application instances. All applications listen to this topic in order to
			be provided the information to setup the ownership of per dialog application instances.
		*/
		Info.Printf("Initializing signalling bus for application %s", app)
		producer := ari.InitProducer(app) // Initialize a new producer channel using the ari.InitProducer function.
		Info.Printf("Starting event handler for application %s", app)
		go runEventHandler(app, producer) // create new websocket connection for every application and pass the producer channel
	}

	go signalCatcher() // listen for os signal to stop the application
	select {}
}

// runEventHandler sets up a websocket connection to an ARI application.
func runEventHandler(s string, producer chan []byte) {
	// Connect to the websocket backend (ARI)
	var ariMessage string
	url := strings.Join([]string{config.Websocket_URL, "?app=", s, "&api_key=", config.WS_User, ":", config.WS_Password}, "")

	Info.Printf("Attempting to connect to ARI websocket at: %s", url)
	ws, err := websocket.Dial(url, "ari", config.Origin)
	if err != nil {
		log.Fatal(err)
	}

	// Start the producer loop. Every message received from the websocket is
	// passed to the PublishMessage() function.
	Info.Printf("Starting producer loop for application %s", s)
	for {
		err = websocket.Message.Receive(ws, &ariMessage) // accept the message from the websocket
		if err != nil {
			log.Fatal(err)
		}
		go PublishMessage(ariMessage, producer) // publish message to the producer channel
	}
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

	switch {
	case info.Type == "StasisStart":
		// since we're starting a new application instance, create the proxy side
		dialogID := ari.UUID()
		Info.Println("New StasisStart found. Created new dialogID of ", dialogID)
		as, err := json.Marshal(ari.AppStart{Application: info.Application, DialogID: dialogID, ServerID: config.ServerID})
		producer <- as

		// TODO: this sleep is required to allow the application time to spin up. In the future we likely want
		// to implement some sort of feedback mechanism in order to remove this sleep timer.
		time.Sleep(50 * time.Millisecond)
		if err != nil {
			return
		}

		Info.Printf("Created new proxy instance mapping for dialog '%s' and channel '%s'", dialogID, info.Channel.ID)
		pi = NewProxyInstance(dialogID)         // create new proxy instance for the dialog
		proxyInstances.Add(info.Channel.ID, pi) // add the dialog to the proxyInstances map to track its life
		exists = true

	case info.Type == "StasisEnd":
		Info.Printf("Ending application instance for channel '%s'", info.Channel.ID)
		// on application end, perform clean up checks
		pi, exists = proxyInstances.Get(info.Channel.ID)
		if exists {
			pi.removeAllObjects()
		}

	case info.Type == "BridgeDestroyed":
		pi, exists = proxyInstances.Get(info.Bridge.ID)
		if exists {
			pi.removeObject(info.Bridge.ID)
		}

	case info.Type == "ChannelDestroyed":
		pi, exists = proxyInstances.Get(info.Channel.ID)
		if exists {
			pi.removeObject(info.Channel.ID)
		}

	// check if prefix is part of the minChan{} struct
	case strings.HasPrefix(info.Type, "Channel"):
		pi, exists = proxyInstances.Get(info.Channel.ID)

	// check if prefix is part of the minBridge{} struct
	case strings.HasPrefix(message.Type, "Bridge"):
		pi, exists = proxyInstances.Get(info.Bridge.ID)

	// check if prefix is part of the minPlay{} struct
	case strings.HasPrefix(message.Type, "Playback"):
		pi, exists = proxyInstances.Get(info.Playback.ID)

	// check if prefix is part of the minRec{} struct (this one uses a Name instead of ID for some reason)
	case strings.HasPrefix(message.Type, "Recording"):
		pi, exists = proxyInstances.Get(info.Recording.Name)

	default:
		Warning.Println("No handler for event type")
		//pi, exists = proxyInstances[]
		// if not matching, then we need to perform checks against the
		// existing map to determine where to send this ARI message.
	}

	// marshal the message back into a string
	busMessage, err := json.Marshal(message)
	if err != nil {
		Error.Println(err)
		return
	}
	Debug.Printf("Bus Data:\n%s\n", busMessage)

	// push the busMessage onto the producer channel
	if exists {
		pi.Events <- busMessage
	}
}

// Add inserts a proxy instance into the global map of active proxy instances.
func (p *proxyInstanceMap) Add(id string, pi *proxyInstance) {
	p.mapLock.Lock()
	defer p.mapLock.Unlock()
	p.instanceMap[id] = pi
}

// Get returns a proxy instance from the global map of active proxy instances.
// returns nil if not found.
func (p *proxyInstanceMap) Get(id string) (*proxyInstance, bool) {
	p.mapLock.RLock()
	defer p.mapLock.RUnlock()
	pi, ok := p.instanceMap[id]
	if !ok {
		return nil, false
	}
	return pi, true
}

// Remove deletes an entry in the global proxyInstanceMap
func (p *proxyInstanceMap) Remove(id string) {
	p.mapLock.Lock()
	defer p.mapLock.Unlock()
	delete(p.instanceMap, id)
}

// shutDown closes the quit channel to signal all of a ProxyInstance's goroutines
// to return
func (p *proxyInstance) shutDown() {
	close(p.quit)
}

// addObject adds an object reference to the proxyInstance mapping
func (p *proxyInstance) addObject(id string) {
	for i := range p.ariObjects {
		if p.ariObjects[i] == id {
			//object already is associated with this proxyInstance
			return
		}
	}
	p.ariObjects = append(p.ariObjects, id)
	proxyInstances.Add(id, p)

}

// removeObject removes an object reference from the proxyInstance mapping
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
	proxyInstances.Remove(id)

	// if there are no more objects, shut'rdown
	if len(p.ariObjects) == 0 {
		p.shutDown()
	}
}

// removeAllObjects will remove all object references from the proxyInstance mapping
func (p *proxyInstance) removeAllObjects() {
	// remove all objects from the map as our application is shutting down.
	for _, obj := range p.ariObjects {
		proxyInstances.Remove(obj)
	}
	p.shutDown() // destroy the application / proxy instance
}

// runCommandConsumer starts the consumer for accepting Commands from
// applications.
func (p *proxyInstance) runCommandConsumer(dialogID string) {
	commandTopic := strings.Join([]string{"commands", dialogID}, "_")
	responseTopic := strings.Join([]string{"responses", dialogID}, "_")
	Debug.Println("Topics are:", commandTopic, " ", responseTopic)
	p.responseChannel = ari.InitProducer(responseTopic)

	// waits for the TopicExists function to return a channel
	select {
	case <-ari.TopicExists(commandTopic):
		p.commandChannel = ari.InitConsumer(commandTopic)
	case <-time.After(10 * time.Second):
		// if the application instance hasn't come up after a period of time, gracefully end the proxy instance
		p.removeAllObjects()
		return
	}

	for {
		select {
		case jsonCommand := <-p.commandChannel:
			go p.processCommand(jsonCommand, p.responseChannel)
		case <-p.quit:
			return
		}
	}
}

// processCommand processes commands from applications and submits them to the
// REST interface.
func (p *proxyInstance) processCommand(jsonCommand []byte, responseProducer chan []byte) {
	var c ari.Command
	var r ari.CommandResponse
	i := ID{ID: "", Name: ""}
	Debug.Printf("jsonCommand is %s\n", string(jsonCommand))
	json.Unmarshal(jsonCommand, &c)
	fullURL := strings.Join([]string{config.Stasis_URL, c.URL, "?api_key=", config.WS_User, ":", config.WS_Password}, "")
	Debug.Printf("fullURL is %s\n", fullURL)
	req, err := http.NewRequest(c.Method, fullURL, bytes.NewBufferString(c.Body))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	Debug.Printf("Response body is %s\n", buf.String())
	json.Unmarshal(buf.Bytes(), &i)
	if i.ID != "" {
		p.addObject(i.ID)
	} else if i.Name != "" {
		p.addObject(i.Name)
	}
	r.ResponseBody = buf.String()
	r.StatusCode = res.StatusCode
	sendJSON, err := json.Marshal(r)
	if err != nil {
		Error.Println(err)
	}
	Debug.Printf("sendJSON is %s\n", string(sendJSON))
	responseProducer <- sendJSON
}
