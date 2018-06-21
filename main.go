package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

// Config is the global app config
type Config struct {
	MsgText    string `default:"Hey Hey, there is a new Turd online"`
	WebHookURL string `required:"true"`
}

// new global config object
var config Config

// SlackMsg implements the SlackAPI for incomming hooks
type SlackMsg struct {
	Text string `json:"text"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	var msg SlackMsg
	msg.Text = config.MsgText
	msgJSON, _ := json.Marshal(msg)
	resp, err := http.Post(config.WebHookURL, "Content-Type: application/json", bytes.NewReader(msgJSON))
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Printf("%s\n", resp.Status)
}

func main() {
	// parse ENV vars
	err := envconfig.Process("notifier", &config)
	if err != nil {
		log.Fatalf("Error in config: %s", err)
	}

	// Server
	http.HandleFunc("/", handler)
	log.Printf("WebhookUrl is: %s", config.WebHookURL)
	log.Printf("%s", "Serving Endpoint '/' on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
