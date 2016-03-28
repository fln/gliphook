package gliphook

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/Sirupsen/logrus"
)

type logrusHook struct {
	url      string
	async    bool
	minLevel logrus.Level
}

// LogrusHook returns a logrus hook that reports error, fatal and panic
// level log messages to Glip via WebHook integration. If async is true
// reporting will be executed in a seperate goroutine. In async mode
// notification errors are not reported back to logrus.
func LogrusHook(url string, async bool, minLevel logrus.Level) logrus.Hook {
	return &logrusHook{
		url:      url,
		async:    async,
		minLevel: minLevel,
	}
}

// Fire implements Hook interface from a github.com/Sirupsen/logrus package.
func (hook *logrusHook) Fire(entry *logrus.Entry) error {
	var bodyBuf bytes.Buffer

	keyOrder := make([]string, 0, len(entry.Data))
	for key := range entry.Data {
		keyOrder = append(keyOrder, key)
	}
	sort.Strings(keyOrder)

	for _, key := range keyOrder {
		fmt.Fprintf(&bodyBuf, "**%s**: %+v\n", key, entry.Data[key])
	}
	event := Notification{
		Activity: fmt.Sprintf("Logrus %v", entry.Level),
		Title:    entry.Message,
		Body:     bodyBuf.String(),
	}
	if hook.async {
		go event.Post(hook.url)
		return nil
	}
	return event.Post(hook.url)
}

// Levels implements Hook interface from a github.com/Sirupsen/logrus package.
func (hook *logrusHook) Levels() []logrus.Level {
	levels := []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
	for i, level := range levels {
		if level == hook.minLevel {
			return levels[i:]
		}
	}
	return []logrus.Level{}
}
