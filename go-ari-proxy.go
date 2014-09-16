package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
)

// Struct Config holds the configuration for the proxy.
// The Config struct contains the information the was unmarshaled from the
// configuration file for ths proxy.
type Config struct {
	Origin        string      `json:"origin"`		// connection to ARI events
	ServerID      string      `json:"server_id"`	// unique server ident
	Applications  []string    `json:"applications"`	// slice of applications to listen for
	Websocket_URL string      `json:"websocket_url"`// websocket to connect to
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

// CreateWS sets up a websocket connection to an ARI application.
func CreateWS(s string, producer chan []byte) {
	// Connect to the websocket backend (ARI)
	var ariMessage string
	url := strings.Join([]string{config.Websocket_URL, "?app=", s, "&api_key=", config.WS_User, ":", config.WS_Password}, "")
	ws, err := websocket.Dial(url, "ari", config.Origin)
	if err != nil {
		log.Fatal(err)
	}

	// Start the producer loop. Every message received from the websocket is passed to the PublishMessage() function.
	for {
		err = websocket.Message.Receive(ws, &ariMessage)		// accept the message from the websocket
		if err != nil {
			log.Fatal(err)
		}
		go PublishMessage(ariMessage, producer)	// publish message to the producer channel
	}
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
	// Setup a new producer for every application we've configured in the configuration file.
	for _, app := range config.Applications {
		producer := ari.InitProducer(config.MessageBus, config.BusConfig, app)	// Initialize a new producer channel using the ari.IntProducer function.
		go CreateWS(app, producer)	// create new websocket connection for every application and pass the producer channel
	}

	go signalCatcher()	// listen for os signal to stop the application
	select {}
}
