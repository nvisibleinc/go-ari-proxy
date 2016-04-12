package ari

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// global variables
var bus MessageBus // global var that holds the current MessageBus interface.

// MessageBus interface contains methods for interacting with the abstracted message bus.
type MessageBus interface {
	InitBus(config interface{}) error
	StartProducer(topic string) (chan []byte, error)
	StartConsumer(topic string) (chan []byte, error)
	TopicExists(topic string) bool
}

// AppInstanceHandler when you start a new App, you pass in a function of type AppInstanceHandler.
// The entry point of the execution of an application instance.
type AppInstanceHandler func(*AppInstance)

// App struct contains information about an ARI application.
// The top level that signals the application instance creation.
type App struct {
	name   string
	Events chan []byte
	Stop   chan bool
}

// AppInstance struct contains the channels necessary for communication to/from
// the various message bus topics and the event channel.
type AppInstance struct {
	commandChannel  chan []byte
	responseChannel chan *CommandResponse
	quit            chan int
	Events          chan *Event
}

// Event struct contains the events we pull off the websocket connection.
type Event struct {
	ServerID  string    `json:"server_id"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	ARI_Body  string    `json:"ari_body"`
}

// AppStart struct contains the initial information for the start of a new application instance.
type AppStart struct {
	Application string `json:"application"`
	DialogID    string `json:"dialog_id"`
	ServerID    string `json:"server_id"`
}

// Command struct contains the command we're passing back to ARI.
type Command struct {
	UniqueID string `json:"unique_id"`
	URL      string `json:"url"`
	Method   string `json:"method"`
	Body     string `json:"body"`
}

// CommandResponse struct contains the response to a Command
type CommandResponse struct {
	UniqueID     string `json:"unique_id"`
	StatusCode   int    `json:"status_code"`
	ResponseBody string `json:"response_body"`
}

// InitLogger is a wrapper function to provide a sane interface to logging messages.
func InitLogger(handle io.Writer, prefix string) *log.Logger {
	return log.New(handle, strings.Join([]string{prefix, ": "}, ""), log.Ldate|log.Ltime|log.Lshortfile)
}

// UUID generates and returns a universally unique identifier.
// TODO(Brad): Replace this with an imported package.
func UUID() string {
	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

// InitBus will initialize a new message bus.
// Abstracts the message bus initialization based on the configuration in order
// to allow the creation of a proxy or client using message bus agnostic
// methods.
func InitBus(busType string, config interface{}) error {
	switch busType {
	case "NATS":
		// Start NATS
		bus = new(NATS)
	case "OSLO":
		// Start an OSLO producer
		log.Fatal("OSLO message bus producer is not yet implemented.")
	case "RABBITMQ":
		// Start a RabbitMQ producer
		bus = new(RabbitMQ)
	default:
		log.Fatal("No bus type was specified for the producer that we recognize.")
	}
	bus.InitBus(config)
	return nil
}

// TopicExists abstracts the basic function provided by the MessageBus interface.
// Spawns a goroutine which loops through and waits for a topic to actually exist.
// Returns a channel immediately which is read by the user of this function to
// determine topic existence or timeout by way of the normal time.After pattern
// in a select{}.
func TopicExists(topic string) <-chan bool {
	c := make(chan bool)
	go func(topic string, c chan bool) {
		for i := 0; i < 20; i++ {
			if bus.TopicExists(topic) {
				c <- true
			}
			time.Sleep(100 * time.Millisecond)
		}
	}(topic, c)
	return c
}

// NewApp creates a new signalling channel for use by an application.
func NewApp() *App {
	var a App
	a.Stop = make(chan bool)
	return &a
}

// Init spawns the goroutine that listens for messages on the signalling channel.
// Creates a new application instance for the client to utilize.
// Passes the AppInstance to the AppInstanceHandler function.
func (a *App) Init(app string, handler AppInstanceHandler) {
	a.Events = InitConsumer(app)
	go func(app string, a *App) {
		for event := range a.Events {
			var as AppStart
			json.Unmarshal(event, &as)
			if as.Application == app {
				ai := new(AppInstance)
				ai.InitAppInstance(as.DialogID)
				go handler(ai)
			}
		}
	}(app, a)
}

// NewAppInstance function is a constructor to allocate the memory of AppInstance.
func NewAppInstance() *AppInstance {
	var a AppInstance
	return &a
}

// InitAppInstance initializes the set of resources necessary for a new application instance.
func (a *AppInstance) InitAppInstance(instanceID string) {
	var err error
	a.Events = make(chan *Event)
	a.responseChannel = make(chan *CommandResponse)
	commandTopic := strings.Join([]string{"commands", instanceID}, "_")
	fmt.Println("Command topic is: ", commandTopic)
	responseTopic := strings.Join([]string{"responses", instanceID}, "_")
	a.commandChannel, err = bus.StartProducer(commandTopic)
	a.commandChannel <- []byte("DUMMY")
	if err != nil {
		fmt.Println(err)
	}
	eventBus, err := bus.StartConsumer(strings.Join([]string{"events", instanceID}, "_"))
	if err != nil {
		fmt.Println(err)
	}
	processEvents(eventBus, a.Events)
	responseBus, err := bus.StartConsumer(responseTopic)
	if err != nil {
		fmt.Println(err)
	}
	a.processCommandResponses(responseBus, a.responseChannel)
}

// InitProducer initializes a new message bus producer.
func InitProducer(topic string) chan []byte {
	producer, err := bus.StartProducer(topic)
	if err != nil {
		fmt.Println(err)
	}
	return producer
}

// InitConsumer initializes a new message bus consumer.
func InitConsumer(topic string) chan []byte {
	consumer, err := bus.StartConsumer(topic)
	if err != nil {
		fmt.Println(err)
	}
	return consumer
}

// processEvents pulls messages off the inboundEvents channel.
// Takes the events which were pulled off the bus, converts them to Event, and
// places onto the parsedEvents channel.
func processEvents(inboundEvents chan []byte, parsedEvents chan *Event) {
	go func(inboundEvents chan []byte, parsedEvents chan *Event) {
		for event := range inboundEvents {
			var e Event
			json.Unmarshal(event, &e)
			parsedEvents <- &e
		}
	}(inboundEvents, parsedEvents)
}

// processCommand is executing the remote command.
// Performs the work of marshaling the command, sending it across the bus, and
// then unmarshaling the data in order to return a command response.
func (a *AppInstance) processCommand(url string, body string, method string) *CommandResponse {
	jsonMessage, err := json.Marshal(Command{URL: url, Method: method, Body: body})
	if err != nil {
		return &CommandResponse{}
	}

	a.commandChannel <- jsonMessage
	for {
		select {
		case r, r_ok := <-a.responseChannel:
			if r_ok {
				return r
			}
		case <-time.After(5 * time.Second):
			return &CommandResponse{}
		}
	}
}

// processCommandResponses is a function for parsing the Command-Response.
// processCommandResponses spawns an anonymous go routine which will listen for
// information on the channel and process them as they arrive.
func (a *AppInstance) processCommandResponses(fromBus chan []byte, toAppInstance chan *CommandResponse) {
	go func(fromBus chan []byte, toAppInstance chan *CommandResponse) {
		for response := range fromBus {
			var cr CommandResponse
			json.Unmarshal(response, &cr)
			toAppInstance <- &cr
		}
	}(fromBus, toAppInstance)
}
