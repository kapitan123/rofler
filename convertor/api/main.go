package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", ConvertVideo)

	// AK TODO fix Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

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
func ConvertVideo(w http.ResponseWriter, r *http.Request) {

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
