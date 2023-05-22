package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// AK TODO we just can't pass all stuff here, we still need an abstraction to group configuration
// temp solution with direct handler function
func setupRouter(r *mux.Router) {
	r.HandleFunc("/pubsub/subscriptions/video-published", convertVideo).Methods("POST")
	r.HandleFunc("/converted-videos/{videoId}", getVideo).Methods("GET")
}

// AK TODO should go to handlerFolder ot shit like this
type ConvertVideoMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		Id   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

// AK TODO need to deserialise this bad boy
type ConvertableVideoPublishedEvent struct {
	VideoUrl          string
	TelegramMessageId int // AK TODO this one is to be able to repost the video after convertion
}

// ConvertVideo receives and processes a Pub/Sub push message with a convertable video
func convertVideo(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var m ConvertVideoMessage
	// byte slice unmarshalling handles base64 decoding.
	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// AK TODO this one is just to test if the message was handled
	name := string(m.Message.Data)
	if name == "" {
		name = "World"
	}
	log.Printf("Hello %s!", name)
}

// ConvertVideo receives and processes a Pub/Sub push message with a convertable video
func getVideo(w http.ResponseWriter, r *http.Request) {
}
