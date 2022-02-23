package app

import (
	"bcm-analyzer/internal/config"
	"bcm-analyzer/internal/httpg"
	"bcm-analyzer/internal/kafkag"
	"bcm-analyzer/internal/logger"

	log "github.com/sirupsen/logrus"
)

var (
	cf  config.Config
	err error
)

func Run() {
	if cf, err = config.New(); err != nil {
		panic(err)
	}

	if err := logger.Run(cf.File); err != nil {
		panic(err)
	}

	kafkag.New(cf.File)

	if err := kafkag.ReadKafkaTopic(); err != nil {
		log.Warnf("kafka ошибка чтения из топика %s", err.Error())
	}

	if err := httpg.HTTPServer(); err != nil {
		panic(err)
	}

}
