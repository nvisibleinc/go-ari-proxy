package ari

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

func buildJSON(params map[string]string) string {
	mapsize := len(params)
	var counter int = 1
	body := bytes.NewBufferString("{")
	for key, value := range params {
		var s string
		if counter < mapsize {
			s = fmt.Sprintf("\"%s\":\"%s\",", key, value)
		} else {
			s = fmt.Sprintf("\"%s\":\"%s\"", key, value)
		}
		body.WriteString(s)
		counter++
	}
	body.WriteString("}")
	return body.String()
}

func (a *AppInstance) ApplicationsList() (*[]Application, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/applications")
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	err = nil
	var r []Application
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) ApplicationsGet(ApplicationName string) (*Application, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/applications/%s", ApplicationName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	switch result.StatusCode {
	case 404:
		err = errors.New("Application does not exist.")
	default:
		err = nil
	}
	var r Application
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) ApplicationsSubscribe(ApplicationName string, EventSource string) (*Application, error) {
	var err error
	paramMap := make(map[string]string)
	paramMap["eventSource"] = EventSource
	url := fmt.Sprintf("/applications/%s/subscription", ApplicationName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 400:
		err = errors.New("Missing parameter.")
	case 404:
		err = errors.New("Application does not exist.")
	case 422:
		err = errors.New("Event source does not exist.")
	default:
		err = nil
	}
	var r Application
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) ApplicationsUnsubscribe(ApplicationName string, EventSource string) (*Application, error) {
	var err error
	paramMap := make(map[string]string)
	paramMap["eventSource"] = EventSource
	url := fmt.Sprintf("/applications/%s/subscription", ApplicationName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "DELETE")
	switch result.StatusCode {
	case 400:
		err = errors.New("Missing parameter; event source scheme not recognized.")
	case 404:
		err = errors.New("Application does not exist.")
	case 409:
		err = errors.New("Application not subscribed to event source.")
	case 422:
		err = errors.New("Event source does not exist.")
	default:
		err = nil
	}
	var r Application
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) AsteriskGetInfo(options ...string) (*AsteriskInfo, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/asterisk/info")
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["only"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	err = nil
	var r AsteriskInfo
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) AsteriskGetGlobalVar(Var string) (*Variable, error) {
	var err error
	paramMap := make(map[string]string)
	paramMap["variable"] = Var
	url := fmt.Sprintf("/asterisk/variable")
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	switch result.StatusCode {
	case 400:
		err = errors.New("Missing variable parameter.")
	default:
		err = nil
	}
	var r Variable
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) AsteriskSetGlobalVar(Var string, options ...string) error {
	var err error
	paramMap := make(map[string]string)
	paramMap["variable"] = Var
	url := fmt.Sprintf("/asterisk/variable")
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["value"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 400:
		err = errors.New("Missing variable parameter.")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) BridgesList() (*[]Bridge, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/bridges")
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	err = nil
	var r []Bridge
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) BridgesCreate(options ...string) (*Bridge, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/bridges")
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["type"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["bridgeId"] = value
			}
		case 2:
			if len(value) > 0 {
				paramMap["name"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	err = nil
	var r Bridge
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) BridgesCreate_Or_Update_With_ID(BridgeID string, options ...string) (*Bridge, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/bridges/%s", BridgeID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["type"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["name"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	err = nil
	var r Bridge
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) BridgesGet(BridgeID string) (*Bridge, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/bridges/%s", BridgeID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	switch result.StatusCode {
	case 404:
		err = errors.New("Bridge not found")
	default:
		err = nil
	}
	var r Bridge
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) BridgesDestroy(BridgeID string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/bridges/%s", BridgeID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "DELETE")
	switch result.StatusCode {
	case 404:
		err = errors.New("Bridge not found")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) BridgesAddChannel(BridgeID string, Channel string, options ...string) error {
	var err error
	paramMap := make(map[string]string)
	paramMap["channel"] = Channel
	url := fmt.Sprintf("/bridges/%s/addChannel", BridgeID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["role"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 400:
		err = errors.New("Channel not found")
	case 404:
		err = errors.New("Bridge not found")
	case 409:
		err = errors.New("Bridge not in Stasis application; Channel currently recording")
	case 422:
		err = errors.New("Channel not in Stasis application")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) BridgesRemoveChannel(BridgeID string, Channel string) error {
	var err error
	paramMap := make(map[string]string)
	paramMap["channel"] = Channel
	url := fmt.Sprintf("/bridges/%s/removeChannel", BridgeID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 400:
		err = errors.New("Channel not found")
	case 404:
		err = errors.New("Bridge not found")
	case 409:
		err = errors.New("Bridge not in Stasis application")
	case 422:
		err = errors.New("Channel not in this bridge")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) BridgesStartMoh(BridgeID string, options ...string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/bridges/%s/moh", BridgeID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["mohClass"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 404:
		err = errors.New("Bridge not found")
	case 409:
		err = errors.New("Bridge not in Stasis application")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) BridgesStopMoh(BridgeID string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/bridges/%s/moh", BridgeID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "DELETE")
	switch result.StatusCode {
	case 404:
		err = errors.New("Bridge not found")
	case 409:
		err = errors.New("Bridge not in Stasis application")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) BridgesPlay(BridgeID string, Media string, options ...string) (*Playback, error) {
	var err error
	paramMap := make(map[string]string)
	paramMap["media"] = Media
	url := fmt.Sprintf("/bridges/%s/play", BridgeID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["lang"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["offsetms"] = value
			}
		case 2:
			if len(value) > 0 {
				paramMap["skipms"] = value
			}
		case 3:
			if len(value) > 0 {
				paramMap["playbackId"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 404:
		err = errors.New("Bridge not found")
	case 409:
		err = errors.New("Bridge not in a Stasis application")
	default:
		err = nil
	}
	var r Playback
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) BridgesPlayWithID(BridgeID string, PlaybackID string, Media string, options ...string) (*Playback, error) {
	var err error
	paramMap := make(map[string]string)
	paramMap["media"] = Media
	url := fmt.Sprintf("/bridges/%s/play/%s", BridgeID, PlaybackID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["lang"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["offsetms"] = value
			}
		case 2:
			if len(value) > 0 {
				paramMap["skipms"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 404:
		err = errors.New("Bridge not found")
	case 409:
		err = errors.New("Bridge not in a Stasis application")
	default:
		err = nil
	}
	var r Playback
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) BridgesRecord(BridgeID string, Name string, Format string, options ...string) (*LiveRecording, error) {
	var err error
	paramMap := make(map[string]string)
	paramMap["name"] = Name
	paramMap["format"] = Format
	url := fmt.Sprintf("/bridges/%s/record", BridgeID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["maxDurationSeconds"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["maxSilenceSeconds"] = value
			}
		case 2:
			if len(value) > 0 {
				paramMap["ifExists"] = value
			}
		case 3:
			if len(value) > 0 {
				paramMap["beep"] = value
			}
		case 4:
			if len(value) > 0 {
				paramMap["terminateOn"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 400:
		err = errors.New("Invalid parameters")
	case 404:
		err = errors.New("Bridge not found")
	case 409:
		err = errors.New("Bridge is not in a Stasis application; A recording with the same name already exists on the system and can not be overwritten because it is in progress or ifExists=fail")
	case 422:
		err = errors.New("The format specified is unknown on this system")
	default:
		err = nil
	}
	var r LiveRecording
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) ChannelsList() (*[]Channel, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/channels")
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	err = nil
	var r []Channel
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) ChannelsOriginate(Endpoint string, options ...string) (*Channel, error) {
	var err error
	paramMap := make(map[string]string)
	paramMap["endpoint"] = Endpoint
	url := fmt.Sprintf("/channels")
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["extension"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["context"] = value
			}
		case 2:
			if len(value) > 0 {
				paramMap["priority"] = value
			}
		case 3:
			if len(value) > 0 {
				paramMap["app"] = value
			}
		case 4:
			if len(value) > 0 {
				paramMap["appArgs"] = value
			}
		case 5:
			if len(value) > 0 {
				paramMap["callerId"] = value
			}
		case 6:
			if len(value) > 0 {
				paramMap["timeout"] = value
			}
		case 7:
			if len(value) > 0 {
				paramMap["variables"] = value
			}
		case 8:
			if len(value) > 0 {
				paramMap["channelId"] = value
			}
		case 9:
			if len(value) > 0 {
				paramMap["otherChannelId"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 400:
		err = errors.New("Invalid parameters for originating a channel.")
	default:
		err = nil
	}
	var r Channel
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) ChannelsGet(ChannelID string) (*Channel, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/channels/%s", ChannelID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	switch result.StatusCode {
	case 404:
		err = errors.New("Channel not found")
	default:
		err = nil
	}
	var r Channel
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) ChannelsOriginateWithID(ChannelID string, Endpoint string, options ...string) (*Channel, error) {
	var err error
	paramMap := make(map[string]string)
	paramMap["endpoint"] = Endpoint
	url := fmt.Sprintf("/channels/%s", ChannelID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["extension"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["context"] = value
			}
		case 2:
			if len(value) > 0 {
				paramMap["priority"] = value
			}
		case 3:
			if len(value) > 0 {
				paramMap["app"] = value
			}
		case 4:
			if len(value) > 0 {
				paramMap["appArgs"] = value
			}
		case 5:
			if len(value) > 0 {
				paramMap["callerId"] = value
			}
		case 6:
			if len(value) > 0 {
				paramMap["timeout"] = value
			}
		case 7:
			if len(value) > 0 {
				paramMap["variables"] = value
			}
		case 8:
			if len(value) > 0 {
				paramMap["otherChannelId"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 400:
		err = errors.New("Invalid parameters for originating a channel.")
	default:
		err = nil
	}
	var r Channel
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) ChannelsHangup(ChannelID string, options ...string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/channels/%s", ChannelID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["reason"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "DELETE")
	switch result.StatusCode {
	case 400:
		err = errors.New("Invalid reason for hangup provided")
	case 404:
		err = errors.New("Channel not found")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) ChannelsContinueInDialplan(ChannelID string, options ...string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/channels/%s/continue", ChannelID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["context"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["extension"] = value
			}
		case 2:
			if len(value) > 0 {
				paramMap["priority"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel not in a Stasis application")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) ChannelsAnswer(ChannelID string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/channels/%s/answer", ChannelID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel not in a Stasis application")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) ChannelsRing(ChannelID string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/channels/%s/ring", ChannelID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel not in a Stasis application")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) ChannelsRingStop(ChannelID string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/channels/%s/ring", ChannelID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "DELETE")
	switch result.StatusCode {
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel not in a Stasis application")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) ChannelsSendDTMF(ChannelID string, options ...string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/channels/%s/dtmf", ChannelID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["dtmf"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["before"] = value
			}
		case 2:
			if len(value) > 0 {
				paramMap["between"] = value
			}
		case 3:
			if len(value) > 0 {
				paramMap["duration"] = value
			}
		case 4:
			if len(value) > 0 {
				paramMap["after"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 400:
		err = errors.New("DTMF is required")
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel not in a Stasis application")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) ChannelsMute(ChannelID string, options ...string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/channels/%s/mute", ChannelID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["direction"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel not in a Stasis application")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) ChannelsUnmute(ChannelID string, options ...string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/channels/%s/mute", ChannelID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["direction"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "DELETE")
	switch result.StatusCode {
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel not in a Stasis application")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) ChannelsHold(ChannelID string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/channels/%s/hold", ChannelID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel not in a Stasis application")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) ChannelsUnhold(ChannelID string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/channels/%s/hold", ChannelID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "DELETE")
	switch result.StatusCode {
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel not in a Stasis application")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) ChannelsStartMoh(ChannelID string, options ...string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/channels/%s/moh", ChannelID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["mohClass"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel not in a Stasis application")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) ChannelsStopMoh(ChannelID string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/channels/%s/moh", ChannelID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "DELETE")
	switch result.StatusCode {
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel not in a Stasis application")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) ChannelsStartSilence(ChannelID string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/channels/%s/silence", ChannelID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel not in a Stasis application")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) ChannelsStopSilence(ChannelID string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/channels/%s/silence", ChannelID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "DELETE")
	switch result.StatusCode {
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel not in a Stasis application")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) ChannelsPlay(ChannelID string, Media string, options ...string) (*Playback, error) {
	var err error
	paramMap := make(map[string]string)
	paramMap["media"] = Media
	url := fmt.Sprintf("/channels/%s/play", ChannelID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["lang"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["offsetms"] = value
			}
		case 2:
			if len(value) > 0 {
				paramMap["skipms"] = value
			}
		case 3:
			if len(value) > 0 {
				paramMap["playbackId"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel not in a Stasis application")
	default:
		err = nil
	}
	var r Playback
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) ChannelsPlayWithID(ChannelID string, PlaybackID string, Media string, options ...string) (*Playback, error) {
	var err error
	paramMap := make(map[string]string)
	paramMap["media"] = Media
	url := fmt.Sprintf("/channels/%s/play/%s", ChannelID, PlaybackID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["lang"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["offsetms"] = value
			}
		case 2:
			if len(value) > 0 {
				paramMap["skipms"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel not in a Stasis application")
	default:
		err = nil
	}
	var r Playback
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) ChannelsRecord(ChannelID string, Name string, Format string, options ...string) (*LiveRecording, error) {
	var err error
	paramMap := make(map[string]string)
	paramMap["name"] = Name
	paramMap["format"] = Format
	url := fmt.Sprintf("/channels/%s/record", ChannelID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["maxDurationSeconds"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["maxSilenceSeconds"] = value
			}
		case 2:
			if len(value) > 0 {
				paramMap["ifExists"] = value
			}
		case 3:
			if len(value) > 0 {
				paramMap["beep"] = value
			}
		case 4:
			if len(value) > 0 {
				paramMap["terminateOn"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 400:
		err = errors.New("Invalid parameters")
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel is not in a Stasis application; the channel is currently bridged with other hcannels; A recording with the same name already exists on the system and can not be overwritten because it is in progress or ifExists=fail")
	case 422:
		err = errors.New("The format specified is unknown on this system")
	default:
		err = nil
	}
	var r LiveRecording
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) ChannelsGetChannelVar(ChannelID string, Var string) (*Variable, error) {
	var err error
	paramMap := make(map[string]string)
	paramMap["variable"] = Var
	url := fmt.Sprintf("/channels/%s/variable", ChannelID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	switch result.StatusCode {
	case 400:
		err = errors.New("Missing variable parameter.")
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel not in a Stasis application")
	default:
		err = nil
	}
	var r Variable
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) ChannelsSetChannelVar(ChannelID string, Var string, options ...string) error {
	var err error
	paramMap := make(map[string]string)
	paramMap["variable"] = Var
	url := fmt.Sprintf("/channels/%s/variable", ChannelID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["value"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 400:
		err = errors.New("Missing variable parameter.")
	case 404:
		err = errors.New("Channel not found")
	case 409:
		err = errors.New("Channel not in a Stasis application")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) ChannelsSnoopChannel(ChannelID string, App string, options ...string) (*Channel, error) {
	var err error
	paramMap := make(map[string]string)
	paramMap["app"] = App
	url := fmt.Sprintf("/channels/%s/snoop", ChannelID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["spy"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["whisper"] = value
			}
		case 2:
			if len(value) > 0 {
				paramMap["appArgs"] = value
			}
		case 3:
			if len(value) > 0 {
				paramMap["snoopId"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 400:
		err = errors.New("Invalid parameters")
	case 404:
		err = errors.New("Channel not found")
	default:
		err = nil
	}
	var r Channel
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) ChannelsSnoopChannelWithID(ChannelID string, SnoopID string, App string, options ...string) (*Channel, error) {
	var err error
	paramMap := make(map[string]string)
	paramMap["app"] = App
	url := fmt.Sprintf("/channels/%s/snoop/%s", ChannelID, SnoopID)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["spy"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["whisper"] = value
			}
		case 2:
			if len(value) > 0 {
				paramMap["appArgs"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 400:
		err = errors.New("Invalid parameters")
	case 404:
		err = errors.New("Channel not found")
	default:
		err = nil
	}
	var r Channel
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) DeviceStatesList() (*[]DeviceState, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/deviceStates")
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	err = nil
	var r []DeviceState
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) DeviceStatesGet(DeviceName string) (*DeviceState, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/deviceStates/%s", DeviceName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	err = nil
	var r DeviceState
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) DeviceStatesUpdate(DeviceName string, DeviceState string) error {
	var err error
	paramMap := make(map[string]string)
	paramMap["deviceState"] = DeviceState
	url := fmt.Sprintf("/deviceStates/%s", DeviceName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "PUT")
	switch result.StatusCode {
	case 404:
		err = errors.New("Device name is missing")
	case 409:
		err = errors.New("Uncontrolled device specified")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) DeviceStatesDelete(DeviceName string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/deviceStates/%s", DeviceName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "DELETE")
	switch result.StatusCode {
	case 404:
		err = errors.New("Device name is missing")
	case 409:
		err = errors.New("Uncontrolled device specified")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) EndpointsList() (*[]Endpoint, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/endpoints")
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	err = nil
	var r []Endpoint
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) EndpointsSendMessage(To string, From string, options ...string) error {
	var err error
	paramMap := make(map[string]string)
	paramMap["to"] = To
	paramMap["from"] = From
	url := fmt.Sprintf("/endpoints/sendMessage")
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["body"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["variables"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "PUT")
	switch result.StatusCode {
	case 404:
		err = errors.New("Endpoint not found")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) EndpointsListByTech(Tech string) (*[]Endpoint, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/endpoints/%s", Tech)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	switch result.StatusCode {
	case 404:
		err = errors.New("Endpoints not found")
	default:
		err = nil
	}
	var r []Endpoint
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) EndpointsGet(Tech string, Resource string) (*Endpoint, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/endpoints/%s/%s", Tech, Resource)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	switch result.StatusCode {
	case 400:
		err = errors.New("Invalid parameters for sending a message.")
	case 404:
		err = errors.New("Endpoints not found")
	default:
		err = nil
	}
	var r Endpoint
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) EndpointsSendMessageToEndpoint(Tech string, Resource string, From string, options ...string) error {
	var err error
	paramMap := make(map[string]string)
	paramMap["from"] = From
	url := fmt.Sprintf("/endpoints/%s/%s/sendMessage", Tech, Resource)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["body"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["variables"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "PUT")
	switch result.StatusCode {
	case 400:
		err = errors.New("Invalid parameters for sending a message.")
	case 404:
		err = errors.New("Endpoint not found")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) EventsEventWebsocket(App string) (*Message, error) {
	var err error
	paramMap := make(map[string]string)
	paramMap["app"] = App
	url := fmt.Sprintf("/events")
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	err = nil
	var r Message
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) EventsUserEvent(EventName string, Application string, options ...string) error {
	var err error
	paramMap := make(map[string]string)
	paramMap["application"] = Application
	url := fmt.Sprintf("/events/user/%s", EventName)
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["source"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["variables"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 404:
		err = errors.New("Application does not exist.")
	case 422:
		err = errors.New("Event source not found.")
	case 400:
		err = errors.New("Invalid even tsource URI or userevent data.")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) MailboxesList() (*[]Mailbox, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/mailboxes")
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	err = nil
	var r []Mailbox
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) MailboxesGet(MailboxName string) (*Mailbox, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/mailboxes/%s", MailboxName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	switch result.StatusCode {
	case 404:
		err = errors.New("Mailbox not found")
	default:
		err = nil
	}
	var r Mailbox
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) MailboxesUpdate(MailboxName string, OldMessages int, NewMessages int) error {
	var err error
	url := fmt.Sprintf("/mailboxes/%s", MailboxName)
	body := fmt.Sprintf("{\"oldMessages\": %d, \"newMessages\": %d }", OldMessages, NewMessages)
	result := a.processCommand(url, body, "PUT")
	switch result.StatusCode {
	case 404:
		err = errors.New("Mailbox not found")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) MailboxesDelete(MailboxName string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/mailboxes/%s", MailboxName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "DELETE")
	switch result.StatusCode {
	case 404:
		err = errors.New("Mailbox not found")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) PlaybacksGet(PlaybackID string) (*Playback, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/playbacks/%s", PlaybackID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	switch result.StatusCode {
	case 404:
		err = errors.New("The playback cannot be found")
	default:
		err = nil
	}
	var r Playback
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) PlaybacksStop(PlaybackID string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/playbacks/%s", PlaybackID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "DELETE")
	switch result.StatusCode {
	case 404:
		err = errors.New("The playback cannot be found")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) PlaybacksControl(PlaybackID string, Operation string) error {
	var err error
	paramMap := make(map[string]string)
	paramMap["operation"] = Operation
	url := fmt.Sprintf("/playbacks/%s/control", PlaybackID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 400:
		err = errors.New("The provided operation parameter was invalid")
	case 404:
		err = errors.New("The playback cannot be found")
	case 409:
		err = errors.New("The operation cannot be performed in the playback's current state")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) RecordingsListStored() (*[]StoredRecording, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/recordings/stored")
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	err = nil
	var r []StoredRecording
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) RecordingsGetStored(RecordingName string) (*StoredRecording, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/recordings/stored/%s", RecordingName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	switch result.StatusCode {
	case 404:
		err = errors.New("Recording not found")
	default:
		err = nil
	}
	var r StoredRecording
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) RecordingsDeleteStored(RecordingName string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/recordings/stored/%s", RecordingName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "DELETE")
	switch result.StatusCode {
	case 404:
		err = errors.New("Recording not found")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) RecordingsCopyStored(RecordingName string, DestinationRecordingName string) (*StoredRecording, error) {
	var err error
	paramMap := make(map[string]string)
	paramMap["destinationRecordingName"] = DestinationRecordingName
	url := fmt.Sprintf("/recordings/stored/%s/copy", RecordingName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 404:
		err = errors.New("Recording not found")
	case 409:
		err = errors.New("A recording with the same name already exists on the system")
	default:
		err = nil
	}
	var r StoredRecording
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) RecordingsGetLive(RecordingName string) (*LiveRecording, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/recordings/live/%s", RecordingName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	switch result.StatusCode {
	case 404:
		err = errors.New("Recording not found")
	default:
		err = nil
	}
	var r LiveRecording
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) RecordingsCancel(RecordingName string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/recordings/live/%s", RecordingName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "DELETE")
	switch result.StatusCode {
	case 404:
		err = errors.New("Recording not found")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) RecordingsStop(RecordingName string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/recordings/live/%s/stop", RecordingName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 404:
		err = errors.New("Recording not found")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) RecordingsPause(RecordingName string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/recordings/live/%s/pause", RecordingName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 404:
		err = errors.New("Recording not found")
	case 409:
		err = errors.New("Recording not in session")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) RecordingsUnpause(RecordingName string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/recordings/live/%s/pause", RecordingName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "DELETE")
	switch result.StatusCode {
	case 404:
		err = errors.New("Recording not found")
	case 409:
		err = errors.New("Recording not in session")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) RecordingsMute(RecordingName string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/recordings/live/%s/mute", RecordingName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "POST")
	switch result.StatusCode {
	case 404:
		err = errors.New("Recording not found")
	case 409:
		err = errors.New("Recording not in session")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) RecordingsUnmute(RecordingName string) error {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/recordings/live/%s/mute", RecordingName)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "DELETE")
	switch result.StatusCode {
	case 404:
		err = errors.New("Recording not found")
	case 409:
		err = errors.New("Recording not in session")
	default:
		err = nil
	}
	return err
}

func (a *AppInstance) SoundsList(options ...string) (*[]Sound, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/sounds")
	for index, value := range options {
		switch index {
		case 0:
			if len(value) > 0 {
				paramMap["lang"] = value
			}
		case 1:
			if len(value) > 0 {
				paramMap["format"] = value
			}
		}
	}
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	err = nil
	var r []Sound
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}

func (a *AppInstance) SoundsGet(SoundID string) (*Sound, error) {
	var err error
	paramMap := make(map[string]string)
	url := fmt.Sprintf("/sounds/%s", SoundID)
	body := buildJSON(paramMap)
	result := a.processCommand(url, body, "GET")
	err = nil
	var r Sound
	json.Unmarshal([]byte(result.ResponseBody), &r)
	return &r, err
}
