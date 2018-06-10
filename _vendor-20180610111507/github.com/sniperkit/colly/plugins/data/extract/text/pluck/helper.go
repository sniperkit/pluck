// +build !log
package plucker

import (
	log "github.com/sirupsen/logrus"
)

// Verbose toggles debug mode
func (p *Plucker) Verbose(makeVerbose bool) {
	if makeVerbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}
