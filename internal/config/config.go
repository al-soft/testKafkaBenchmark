package config

import (
	"flag"
	"os"
)

// Config данные о конфигурационном файле
type Config struct {
	DefaultFile string
	File        *string
}

// New инициализация наименования конфигурационного файла и проверка существования
func New() (Config, error) {
	cf := Config{}
	cf.DefaultFile = "config.local.yml"
	cf.File = flag.String("config", cf.DefaultFile, "path to the config file")
	flag.Parse()

	if _, err := os.Stat(*cf.File); os.IsNotExist(err) {
		return cf, err
	}

	return cf, nil
}
