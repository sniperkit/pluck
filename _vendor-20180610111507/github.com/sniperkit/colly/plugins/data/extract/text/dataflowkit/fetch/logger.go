package fetch

import (
	"github.com/slotix/dataflowkit/logger"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = log.NewLogger(true)
}
