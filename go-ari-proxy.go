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
}

// PublishMessage takes an ARI event from the websocket and places it on the 
// producer channel.
// Accepts two arguments:
// * a string containing the ARI message
// * a producer channel
func PublishMessage(ariMessage string, producer chan []byte) {
	// unmarshal into an ari.Event so we can append some extra information
	var message ari.Event
	json.Unmarshal([]byte(ariMessage), &message)
	message.ServerID = config.ServerID
	message.Timestamp = time.Now()
	message.ARI_Body = ariMessage

	// marshal the message back into a string
	busMessage, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[DEBUG] Bus Data:\n%s", busMessage)

	// please the busMessage onto the producer channel
	producer <- busMessage
}

// ConsumeCommand is used for consuming commands from an Asterisk instance.
func ConsumeCommand() {

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

// runCommandConsumer starts the consumer for accepting Commands from
// applications.
func runCommandConsumer(app string) {
	consumer := ari.InitConsumer(strings.Join([]string{app, "commands"},"_"), config.MessageBus, config.BusConfig)
	responseProducer := ari.InitProducer(strings.Join([]string{app, "responses"},"_"), config.MessageBus, config.BusConfig)
	for jsonCommand := range consumer {
		fmt.Printf("runCommandConsumer: %s\n", jsonCommand)
		go processCommand(jsonCommand, responseProducer)
	}
}

func processCommand(jsonCommand []byte, responseProducer chan []byte) {
	var c ari.Command
	var r ari.CommandResponse
	json.Unmarshal(jsonCommand, &c)
	fullURL := strings.Join([]string{config.Stasis_URL, c.URL, "?api_key=", config.WS_User, ":", config.WS_Password }, "")
	fmt.Println(fullURL)
	req, err := http.NewRequest(c.Method, fullURL, bytes.NewBufferString(c.Body))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	r.ResponseBody = buf.String()
	r.UniqueID = c.UniqueID
	r.StatusCode = res.StatusCode
	sendJSON, err := json.Marshal(r)
	if err !=nil {
		fmt.Println(err)
	}
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
	for _, app := range config.Applications {
		producer := ari.InitProducer(app, config.MessageBus, config.BusConfig)	// Initialize a new producer channel using the ari.IntProducer function.
		go runEventHandler(app, producer)	// create new websocket connection for every application and pass the producer channel
		go runCommandConsumer(app)
	}

	go signalCatcher()	// listen for os signal to stop the application
	select {}
}
