package gliphook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotificationPost(t *testing.T) {
	expected := Notification{
		Icon:     "http://localhost/img.png",
		Activity: "act",
		Title:    "title",
		Body:     "body",
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err, "Error reading request")

		actual := Notification{}
		err = json.Unmarshal(body, &actual)
		assert.Nil(t, err, "Error unmarshaling request")

		assert.Equal(t, expected, actual, "Notification was modified in transit")

		code := http.StatusNoContent
		http.Error(w, http.StatusText(code), code)
	})
	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	// HTTP 204 response from the server
	assert.Nil(t, expected.Post(server.URL+"/ok"), "Valid request triggered an error")

	// HTTP 500 response from the server
	assert.NotNil(t, expected.Post(server.URL+"/error"), "HTTP 500 response does not trigger an error")

	// Trigger HTTP transport error
	assert.NotNil(t, expected.Post("example://invalid/url"), "HTTP transport error does not trigger an errror")
}
