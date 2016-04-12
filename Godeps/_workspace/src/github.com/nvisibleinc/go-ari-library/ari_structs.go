package ari

type Channel struct {
	Id           string      `json:"id"`
	Name         string      `json:"name"`
	State        string      `json:"state"`
	Caller       CallerID    `json:"caller"`
	Connected    CallerID    `json:"connected"`
	Accountcode  string      `json:"accountcode"`
	Dialplan     DialplanCEP `json:"dialplan"`
	Creationtime string      `json:"creationtime"`
}

type BridgeDestroyed struct {
	Bridge      Bridge `json:"bridge"`
	Application string `json:"application"`
	Timestamp   string `json:"timestamp"`
	Type        string `json:"type"`
}

type BridgeCreated struct {
	Bridge      Bridge `json:"bridge"`
	Application string `json:"application"`
	Timestamp   string `json:"timestamp"`
	Type        string `json:"type"`
}

type Playback struct {
	Id         string `json:"id"`
	Media_Uri  string `json:"media_uri"`
	Target_Uri string `json:"target_uri"`
	Language   string `json:"language"`
	State      string `json:"state"`
}

type DeviceState struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

type ChannelDtmfReceived struct {
	Digit       string  `json:"digit"`
	Duration_Ms int     `json:"duration_ms"`
	Channel     Channel `json:"channel"`
	Application string  `json:"application"`
	Timestamp   string  `json:"timestamp"`
	Type        string  `json:"type"`
}

type ChannelHangupRequest struct {
	Cause       int     `json:"cause"`
	Soft        bool    `json:"soft"`
	Channel     Channel `json:"channel"`
	Application string  `json:"application"`
	Timestamp   string  `json:"timestamp"`
	Type        string  `json:"type"`
}

type TextMessageReceived struct {
	Message     TextMessage `json:"message"`
	Endpoint    Endpoint    `json:"endpoint"`
	Application string      `json:"application"`
	Timestamp   string      `json:"timestamp"`
	Type        string      `json:"type"`
}

type LiveRecording struct {
	Name             string `json:"name"`
	Format           string `json:"format"`
	Target_Uri       string `json:"target_uri"`
	State            string `json:"state"`
	Duration         int    `json:"duration"`
	Talking_Duration int    `json:"talking_duration"`
	Silence_Duration int    `json:"silence_duration"`
	Cause            string `json:"cause"`
}

type ChannelEnteredBridge struct {
	Bridge      Bridge  `json:"bridge"`
	Channel     Channel `json:"channel"`
	Application string  `json:"application"`
	Timestamp   string  `json:"timestamp"`
	Type        string  `json:"type"`
}

type Bridge struct {
	Id           string   `json:"id"`
	Technology   string   `json:"technology"`
	Bridge_Type  string   `json:"bridge_type"`
	Bridge_Class string   `json:"bridge_class"`
	Creator      string   `json:"creator"`
	Name         string   `json:"name"`
	Channels     []string `json:"channels"`
}

type BridgeMerged struct {
	Bridge      Bridge `json:"bridge"`
	Bridge_From Bridge `json:"bridge_from"`
	Application string `json:"application"`
	Timestamp   string `json:"timestamp"`
	Type        string `json:"type"`
}

type BuildInfo struct {
	Os      string `json:"os"`
	Kernel  string `json:"kernel"`
	Options string `json:"options"`
	Machine string `json:"machine"`
	Date    string `json:"date"`
	User    string `json:"user"`
}

type SystemInfo struct {
	Version   string `json:"version"`
	Entity_ID string `json:"entity_id"`
}

type ConfigInfo struct {
	Name             string  `json:"name"`
	Default_Language string  `json:"default_language"`
	Max_Channels     int     `json:"max_channels"`
	Max_Open_Files   int     `json:"max_open_files"`
	Max_Load         float64 `json:"max_load"`
	Setid            SetId   `json:"setid"`
}

type RecordingFailed struct {
	Recording   LiveRecording `json:"recording"`
	Application string        `json:"application"`
	Timestamp   string        `json:"timestamp"`
	Type        string        `json:"type"`
}

type PlaybackFinished struct {
	Playback    Playback `json:"playback"`
	Application string   `json:"application"`
	Timestamp   string   `json:"timestamp"`
	Type        string   `json:"type"`
}

type ChannelUserevent struct {
	Eventname   string   `json:"eventname"`
	Channel     Channel  `json:"channel"`
	Bridge      Bridge   `json:"bridge"`
	Endpoint    Endpoint `json:"endpoint"`
	Userevent   string   `json:"userevent"`
	Application string   `json:"application"`
	Timestamp   string   `json:"timestamp"`
	Type        string   `json:"type"`
}

type ChannelCallerId struct {
	Caller_Presentation     int     `json:"caller_presentation"`
	Caller_Presentation_Txt string  `json:"caller_presentation_txt"`
	Channel                 Channel `json:"channel"`
	Application             string  `json:"application"`
	Timestamp               string  `json:"timestamp"`
	Type                    string  `json:"type"`
}

type EndpointStateChange struct {
	Endpoint    Endpoint `json:"endpoint"`
	Application string   `json:"application"`
	Timestamp   string   `json:"timestamp"`
	Type        string   `json:"type"`
}

type BridgeAttendedTransfer struct {
	Result                       string  `json:"result"`
	Transferer_First_Leg         Channel `json:"transferer_first_leg"`
	Transferer_Second_Leg        Channel `json:"transferer_second_leg"`
	Replace_Channel              Channel `json:"replace_channel"`
	Is_External                  bool    `json:"is_external"`
	Transferer_First_Leg_Bridge  Bridge  `json:"transferer_first_leg_bridge"`
	Transferer_Second_Leg_Bridge Bridge  `json:"transferer_second_leg_bridge"`
	Destination_Bridge           string  `json:"destination_bridge"`
	Destination_Link_Second_Leg  Channel `json:"destination_link_second_leg"`
	Destination_Threeway_Channel Channel `json:"destination_threeway_channel"`
	Destination_Threeway_Bridge  Bridge  `json:"destination_threeway_bridge"`
	Transferee                   Channel `json:"transferee"`
	Transfer_Target              Channel `json:"transfer_target"`
	Destination_Type             string  `json:"destination_type"`
	Destination_Application      string  `json:"destination_application"`
	Destination_Link_First_Leg   Channel `json:"destination_link_first_leg"`
	Application                  string  `json:"application"`
	Timestamp                    string  `json:"timestamp"`
	Type                         string  `json:"type"`
}

type Mailbox struct {
	Name         string `json:"name"`
	Old_Messages int    `json:"old_messages"`
	New_Messages int    `json:"new_messages"`
}

type Sound struct {
	Id      string           `json:"id"`
	Text    string           `json:"text"`
	Formats []FormatLangPair `json:"formats"`
}

type Dialed struct {
}

type ChannelStateChange struct {
	Channel     Channel `json:"channel"`
	Application string  `json:"application"`
	Timestamp   string  `json:"timestamp"`
	Type        string  `json:"type"`
}

type ChannelCreated struct {
	Channel     Channel `json:"channel"`
	Application string  `json:"application"`
	Timestamp   string  `json:"timestamp"`
	Type        string  `json:"type"`
}

type ChannelTalkingStarted struct {
	Channel     Channel `json:"channel"`
	Application string  `json:"application"`
	Timestamp   string  `json:"timestamp"`
	Type        string  `json:"type"`
}

type ChannelLeftBridge struct {
	Bridge      Bridge  `json:"bridge"`
	Channel     Channel `json:"channel"`
	Application string  `json:"application"`
	Timestamp   string  `json:"timestamp"`
	Type        string  `json:"type"`
}

type StoredRecording struct {
	Name   string `json:"name"`
	Format string `json:"format"`
}

type FormatLangPair struct {
	Language string `json:"language"`
	Format   string `json:"format"`
}

type StatusInfo struct {
	Startup_Time     string `json:"startup_time"`
	Last_Reload_Time string `json:"last_reload_time"`
}

type Endpoint struct {
	Technology  string   `json:"technology"`
	Resource    string   `json:"resource"`
	State       string   `json:"state"`
	Channel_IDs []string `json:"channel_ids"`
}

type DeviceStateChanged struct {
	Device_State DeviceState `json:"device_state"`
	Application  string      `json:"application"`
	Timestamp    string      `json:"timestamp"`
	Type         string      `json:"type"`
}

type RecordingStarted struct {
	Recording   LiveRecording `json:"recording"`
	Application string        `json:"application"`
	Timestamp   string        `json:"timestamp"`
	Type        string        `json:"type"`
}

type StasisEnd struct {
	Channel     Channel `json:"channel"`
	Application string  `json:"application"`
	Timestamp   string  `json:"timestamp"`
	Type        string  `json:"type"`
}

type Application struct {
	Name         string   `json:"name"`
	Channel_IDs  []string `json:"channel_ids"`
	Bridge_IDs   []string `json:"bridge_ids"`
	Endpoint_IDs []string `json:"endpoint_ids"`
	Device_Names []string `json:"device_names"`
}

type MissingParams struct {
	Params []string `json:"params"`
	Type   string   `json:"type"`
}

type Dial struct {
	Caller      Channel `json:"caller"`
	Peer        Channel `json:"peer"`
	Forward     string  `json:"forward"`
	Forwarded   Channel `json:"forwarded"`
	Dialstring  string  `json:"dialstring"`
	Dialstatus  string  `json:"dialstatus"`
	Application string  `json:"application"`
	Timestamp   string  `json:"timestamp"`
	Type        string  `json:"type"`
}

type AsteriskInfo struct {
	Build  BuildInfo  `json:"build"`
	System SystemInfo `json:"system"`
	Config ConfigInfo `json:"config"`
	Status StatusInfo `json:"status"`
}

type TextMessageVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ChannelTalkingFinished struct {
	Channel     Channel `json:"channel"`
	Duration    int     `json:"duration"`
	Application string  `json:"application"`
	Timestamp   string  `json:"timestamp"`
	Type        string  `json:"type"`
}

type ChannelVarset struct {
	Variable    string  `json:"variable"`
	Value       string  `json:"value"`
	Channel     Channel `json:"channel"`
	Application string  `json:"application"`
	Timestamp   string  `json:"timestamp"`
	Type        string  `json:"type"`
}

type StasisStart struct {
	Args            []string `json:"args"`
	Channel         Channel  `json:"channel"`
	Replace_Channel Channel  `json:"replace_channel"`
	Application     string   `json:"application"`
	Timestamp       string   `json:"timestamp"`
	Type            string   `json:"type"`
}

type ChannelDestroyed struct {
	Cause       int     `json:"cause"`
	Cause_Txt   string  `json:"cause_txt"`
	Channel     Channel `json:"channel"`
	Application string  `json:"application"`
	Timestamp   string  `json:"timestamp"`
	Type        string  `json:"type"`
}

type ChannelDialplan struct {
	Channel           Channel `json:"channel"`
	Dialplan_App      string  `json:"dialplan_app"`
	Dialplan_App_Data string  `json:"dialplan_app_data"`
	Application       string  `json:"application"`
	Timestamp         string  `json:"timestamp"`
	Type              string  `json:"type"`
}

type SetId struct {
	User  string `json:"user"`
	Group string `json:"group"`
}

type CallerID struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}

type PlaybackStarted struct {
	Playback    Playback `json:"playback"`
	Application string   `json:"application"`
	Timestamp   string   `json:"timestamp"`
	Type        string   `json:"type"`
}

type BridgeBlindTransfer struct {
	Channel         Channel `json:"channel"`
	Replace_Channel Channel `json:"replace_channel"`
	Transferee      Channel `json:"transferee"`
	Exten           string  `json:"exten"`
	Context         string  `json:"context"`
	Result          string  `json:"result"`
	Is_External     bool    `json:"is_external"`
	Bridge          Bridge  `json:"bridge"`
	Application     string  `json:"application"`
	Timestamp       string  `json:"timestamp"`
	Type            string  `json:"type"`
}

type Message struct {
	Type string `json:"type"`
}

type RecordingFinished struct {
	Recording   LiveRecording `json:"recording"`
	Application string        `json:"application"`
	Timestamp   string        `json:"timestamp"`
	Type        string        `json:"type"`
}

type DialplanCEP struct {
	Context  string `json:"context"`
	Exten    string `json:"exten"`
	Priority uint64 `json:"priority"`
}

type Variable struct {
	Value string `json:"value"`
}

type TextMessage struct {
	From      string                `json:"from"`
	To        string                `json:"to"`
	Body      string                `json:"body"`
	Variables []TextMessageVariable `json:"variables"`
}

type ApplicationReplaced struct {
	Application string `json:"application"`
	Timestamp   string `json:"timestamp"`
	Type        string `json:"type"`
}
