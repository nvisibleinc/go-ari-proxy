package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"code.google.com/p/go.net/websocket"
	"go-ari-library"
)

var (
	config Config
)

type Config struct {
	Origin        string      `json:"origin"`
	ServerID      string      `json:"server_id"`
	Applications  []string    `json:"applications"`
	Websocket_URL string      `json:"websocket_url"`
	WS_User       string      `json:"ws_user"`
	WS_Password   string      `json:"ws_password"`
	MessageBus    string      `json:"message_bus"`
	BusConfig     interface{} `json:"bus_config"`
}

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

func PublishMessage(ariMessage string, producer chan []byte) {
	var message ari.Event
	json.Unmarshal([]byte(ariMessage), &message)
	message.ServerID = config.ServerID
	message.Timestamp = time.Now()
	message.ARI_Body = ariMessage

	busMessage, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[DEBUG] Bus Data:\n%s", busMessage)
	producer <- busMessage
}

func ConsumeCommand() {

}

func CreateWS(s string, producer chan []byte) {
	// connect to the websocket backend (ARI)
	var ariMessage string
	url := strings.Join([]string{config.Websocket_URL, "?app=", s, "&api_key=", config.WS_User, ":", config.WS_Password}, "")
	ws, err := websocket.Dial(url, "ari", config.Origin)
	if err != nil {
		log.Fatal(err)
	}

	// start the producer loop. Every message received from the websocket backend should be published to the message bus as Stasis_Events
	for {
		err = websocket.Message.Receive(ws, &ariMessage)
		if err != nil {
			log.Fatal(err)
		}
		go PublishMessage(ariMessage, producer)
	}
}

func signalCatcher() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	sig := <-ch
	log.Printf("Signal received: %v", sig)
	os.Exit(0)
}

func main() {

	// p *nsq.Producer
	for _, app := range config.Applications {
		producer := ari.InitProducer(config.MessageBus, config.BusConfig, app)
		go CreateWS(app, producer)
	}

	go signalCatcher()
	select {}
}
