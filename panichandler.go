package gliphook

import (
	"fmt"
	"net/http"
)

// PanicHandler creates a panic handler that dumps panic error and stack
// trace to Glip via WebHook.
func PanicHandler(url string) func(err interface{}, stack []byte) {
	return func(err interface{}, stack []byte) {
		e := Notification{
			Activity: "Panic",
			Title:    fmt.Sprintf("%s", err),
			Body:     string(stack),
		}
		go e.Post(url)
	}
}

// HTTPPanicHandler creates a new panic handler (used for catching panics in
// HTTP handlers). This handler will dump panic error and stack trace to a given
// Glip WebHook URL.
func HTTPPanicHandler(url string) func(r *http.Request, val interface{}, stack []byte) {
	return func(r *http.Request, val interface{}, stack []byte) {
		e := Notification{
			Activity: "Panic",
			Title:    fmt.Sprintf("%s", val),
			Body: fmt.Sprintf("%s %s %s\n---- Stack -----\n%s",
				r.RemoteAddr, r.Method, r.URL.String(),
				stack,
			),
		}
		go e.Post(url)
	}
}
