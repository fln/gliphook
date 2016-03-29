package gliphook

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// Notification represents an entry Glip glip chat that can be sent via Glip
// WebHook integration.
type Notification struct {
	Icon     string `json:"icon,omitempty"`     // absolute URL to an image
	Activity string `json:"activity,omitempty"` // type of the notification, displayed in bold
	Title    string `json:"title,omitempty"`    // first line of the notification
	Body     string `json:"body,omitempty"`     // message, can contain a simple form of markdown syntax
}

// Post sends a notification to Glip via WebHook url.
func (n *Notification) Post(url string) error {
	json, err := json.Marshal(n)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(json))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New(resp.Status)
	}

	return nil
}
