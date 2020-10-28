package utils

import log "github.com/sirupsen/logrus"

func InitLogs(level log.Level) {
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	log.SetLevel(level)
}
