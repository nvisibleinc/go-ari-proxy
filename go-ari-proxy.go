package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
	"os"
  "os/signal"
	"syscall"

	"code.google.com/p/go.net/websocket"
	"github.com/bitly/go-nsq"
	"github.com/bitly/nsq/util"
)

var (
	config Config
)

type Config struct {
	Origin        string `json:"origin"`
	ServerID      string `json:"server_id"`
	Applications  []string `json:"applications"`
	Websocket_URL string `json:"websocket_url"`
	WS_User       string `json:"ws_user"`
	WS_Password   string `json:"ws_password"`
	NSQ_Addr      string `json:"nsq_addr"`
}

type NV_Event struct {
  ServerID    string `json:"server_id"`
  Timestamp   time.Time `json:"timestamp"`
  Type        string `json:"type"`
  ARI_Event   string `json:"ari_event"`
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
}

func PublishMessage(ariMessage,ariApplication string, p *nsq.Producer) {
  var message NV_Event
  json.Unmarshal([]byte(ariMessage), &message)
  message.ServerID = config.ServerID
  message.Timestamp = time.Now()
  message.ARI_Event = ariMessage

	busMessage, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
  fmt.Printf("[DEBUG] Bus Data:\n%s", busMessage)
	p.Publish(ariApplication, []byte(busMessage))
}

func ConsumeCommand() {
  
}

func CreateWS(s string) {
	// connect to the websocket backend (ARI)
	var ariMessage string
	url := strings.Join([]string{config.Websocket_URL, "?app=", s, "&api_key=", config.WS_User, ":", config.WS_Password}, "")
	ws, err := websocket.Dial(url, "ari", config.Origin)
	if err != nil {
		log.Fatal(err)
	}

	// load new nsq instance from configuration file
	nsqcfg := nsq.NewConfig()
	nsqcfg.UserAgent = fmt.Sprintf("to_nsq/%s go-nsq/%s", util.BINARY_VERSION, nsq.VERSION)
	producer, err := nsq.NewProducer(config.NSQ_Addr, nsqcfg)
	if err != nil {
		log.Fatal(err)
	}
	
	// start the producer loop. Every message received from the websocket backend should be published to the message bus as Stasis_Events
	for {
		err = websocket.Message.Receive(ws, &ariMessage)
		if err != nil {
			log.Fatal(err)
		}
		go PublishMessage(ariMessage, s, producer)
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
	for _, app := range config.Applications {
		go CreateWS(app)
	}
	
	go signalCatcher()
    select {}
}