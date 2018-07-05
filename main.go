package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"text/template"

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
const slackMsg = `{
		"attachments": [
			{
				"title": "{{ .Title }}",
				"image_url": "{{ .ImageURL }}"
			},
			{
				"fallback": "View in Admin Dashboard",
				"title": "View in Admin Dashboard",
				"color": "#3AA3E3",
				"attachment_type": "default",
				"actions": [
						{
								"name": "dashboard",
								"text": "Admin Dashboard",
								"type": "button",
								"value": "{{ .DashboardURL }}"
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

	// fmt.Println(string(msg))

	// Post the custom turd msg to slack
	resp, err := http.Post(config.WebHookURL, "Content-Type: application/json", bytes.NewReader(msg))
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
