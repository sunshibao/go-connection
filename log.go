package connection

import (
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()
)

func SetLogger(log *logrus.Logger) {
	logger = log
}
