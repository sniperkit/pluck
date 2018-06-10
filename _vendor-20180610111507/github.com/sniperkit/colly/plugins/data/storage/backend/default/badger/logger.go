package badgercache

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var log = logrus.New()

// type logFields logrus.Fields
type Fields logrus.Fields

// WithFields is an alias for logrus.WithFields.
func LogWithFields(f Fields) *logrus.Entry {
	return logrus.WithFields(logrus.Fields(f))
}

// type logFields map[string]interface{}

func init() {
	log.Formatter = new(prefixed.TextFormatter)
	log.Level = logrus.DebugLevel
}
