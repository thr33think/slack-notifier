package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config is the global app config
type Config struct {
	WebHookURL string `required:"true"`
}

// new global config object
var config Config

// SlackMsg implements the SlackAPI for incomming hooks
const slackMsg = `{
		"attachments": [
			{
				"title": "{{ .Title }}",
				"image_url": "{{ .ImageURL }}"
			},
			{
				"fallback": "Check turd here {{ .DashboardURL }}",
				"actions": [
						{
								"type": "button",
								"text": "Admin Dashboard",
								"url": "{{ .DashboardURL }}"
						}
				]
			}
		]
	}`

// TurdMsg defines a custom turd msg that will be posted to slack
type TurdMsg struct {
	Title        string
	ImageURL     string
	DashboardURL string
}

func newTurdMsg(title, imageURL, dashboardURL string) ([]byte, error) {
	// New turdMsg entity
	var turdMsg = TurdMsg{
		Title:        title,
		ImageURL:     imageURL,
		DashboardURL: dashboardURL}

	// RenderedTurdMsg holds the rendered template
	var renderedTurdMsg bytes.Buffer

	// Render the template
	t := template.Must(template.New("turdMsg").Parse(slackMsg))
	err := t.Execute(&renderedTurdMsg, turdMsg)
	if err != nil {
		return nil, err
	}

	return renderedTurdMsg.Bytes(), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	// parse query values
	title := r.URL.Query().Get("title")
	imageURL := r.URL.Query().Get("imageURL")
	dashboardURL := r.URL.Query().Get("dashboardURL")

	// create a new turd msg based on the query values
	msg, err := newTurdMsg(title, imageURL, dashboardURL)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	// Post the custom turd msg to slack
	go func() {
		client := http.Client{
			Timeout: 2 * time.Second,
		}
		req, err := http.NewRequest("POST", config.WebHookURL, bytes.NewReader(msg))
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		fmt.Printf("%s\n", resp.Status)
	}()
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
