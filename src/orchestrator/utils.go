package main

import (
	"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type logData struct {
	uuid   string
	client mqtt.Client
	name   string
}

const playAudioFileTopic string = "bloob/%s/cores/audio_playback_util/play_file"
const recordSpeechTopic string = "bloob/%s/cores/audio_recorder_util/record_speech"
const transcribeAudioTopic string = "bloob/%s/cores/stt_util/transcribe"
const intentParseTopic string = "bloob/%s/cores/intent_parser_util/run"
const ttsTopic string = "bloob/%s/cores/tts_util/run"
const coreTopic string = "bloob/%s/cores/%s/run"
const instantIntentTopic string = "bloob/%s/instant_intents"

const thinkingTopic string = "bloob/%s/thinking"
const recordingTopic string = "bloob/%s/recording"

func playAudioFile(audio string, uuid string, id string, client mqtt.Client) {
	audioPlaybackMessage := map[string]string{
		"id":    id,
		"audio": audio,
	}
	audioPlaybackJson, err := json.Marshal(audioPlaybackMessage)
	if err != nil {
		bLogFatal(err.Error(), l)

	}

	client.Publish(fmt.Sprintf(playAudioFileTopic, uuid), bloobQOS, false, audioPlaybackJson)
}

func startRecordingAudio(uuid string, id string, client mqtt.Client) {
	audioRecordMessage := map[string]string{
		"id": id,
	}
	audioRecordJson, err := json.Marshal(audioRecordMessage)
	if err != nil {
		bLogFatal(err.Error(), l)
	}
	client.Publish(fmt.Sprintf(recordSpeechTopic, uuid), bloobQOS, false, audioRecordJson)
}

func transcribeAudio(audio string, uuid string, id string, client mqtt.Client) {
	audioTranscribeMessage := map[string]string{
		"id":    id,
		"audio": audio,
	}
	audioTranscribeJson, err := json.Marshal(audioTranscribeMessage)
	if err != nil {
		bLogFatal(err.Error(), l)

	}

	client.Publish(fmt.Sprintf(transcribeAudioTopic, uuid), bloobQOS, false, audioTranscribeJson)
}

func intentParseText(text string, uuid string, id string, client mqtt.Client) {
	intentParseMessage := map[string]string{
		"id":   id,
		"text": text,
	}
	intentParseJson, err := json.Marshal(intentParseMessage)
	if err != nil {
		bLogFatal(err.Error(), l)

	}
	client.Publish(fmt.Sprintf(intentParseTopic, uuid), bloobQOS, false, intentParseJson)
}

func sendIntentToCore(intent string, text string, coreId string, uuid string, id string, client mqtt.Client) {
	coreMessage := map[string]string{
		"id":      id,
		"intent":  intent,
		"core_id": coreId,
		"text":    text,
	}
	coreMessageJson, err := json.Marshal(coreMessage)
	if err != nil {
		bLogFatal(err.Error(), l)
	}
	client.Publish(fmt.Sprintf(coreTopic, uuid, coreId), bloobQOS, false, coreMessageJson)
}

func speakText(text string, uuid string, id string, client mqtt.Client) {
	ttsMessage := map[string]string{
		"id":   id,
		"text": text,
	}
	ttsMessageJson, err := json.Marshal(ttsMessage)
	if err != nil {
		bLogFatal(err.Error(), l)
	}
	client.Publish(fmt.Sprintf(ttsTopic, uuid), bloobQOS, false, ttsMessageJson)
}

func setThinking(state bool, uuid string, client mqtt.Client) {
	thinkMessage := map[string]bool{
		"is_thinking": state,
	}
	thinkMessageJson, err := json.Marshal(thinkMessage)
	if err != nil {
		bLogFatal(err.Error(), l)
	}
	client.Publish(fmt.Sprintf(thinkingTopic, uuid), bloobQOS, true, thinkMessageJson)
}

func setRecording(state bool, uuid string, client mqtt.Client) {
	recordingMessage := map[string]bool{
		"is_recording": state,
	}
	recordingMessageJson, err := json.Marshal(recordingMessage)
	if err != nil {
		bLogFatal(err.Error(), l)
	}
	client.Publish(fmt.Sprintf(recordingTopic, uuid), bloobQOS, true, recordingMessageJson)
}

func bLog(text string, ld logData) {
	logMessage := fmt.Sprintf("[%s] %s", ld.name, text)
	if ld.client != nil && ld.uuid != "" {
		ld.client.Publish(fmt.Sprintf("bloob/%s/logs", ld.uuid), bloobQOS, false, logMessage)
	} else {
		logMessage = fmt.Sprintf("[NO MQTT LOGS] %s", logMessage)
	}

	log.Print(logMessage)
}

func bLogFatal(text string, ld logData) {
	logMessage := fmt.Sprintf("!FATAL! [%s] %s", ld.name, text)

	if ld.client != nil {
		ld.client.Publish(fmt.Sprintf("bloob/%s/logs", ld.uuid), bloobQOS, false, logMessage)
	}
	log.Fatal(logMessage)
}