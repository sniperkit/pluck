package sitemap

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var log *logrus.Logger

func InitLogger() {
	log = logrus.New()
	log.Formatter = new(prefixed.TextFormatter)
	log.Level = logrus.DebugLevel
}

// String return the string format of the sitemap
func (s *SitemapCollector) AttachLogger(logger *logrus.Logger) {
	if log == nil {
		log = logrus.New()
		log.Formatter = new(prefixed.TextFormatter)
		log.Level = logrus.DebugLevel
	}
}
