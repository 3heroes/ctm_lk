package config

import (
	"ctm_lk/pkg/logger"
	"encoding/json"
	"os"
)

type ConfigJson struct {
	ServerHost      string
	ServerHttpPort  string
	ServerHttpsPort string
	Db              DB
}

type DB struct {
	Server   string
	Port     int
	User     string
	Password string
	Database string
}

func defaultJsonBytes() []byte {
	cj := ConfigJson{
		ServerHost:      "localhost",
		ServerHttpPort:  "80",
		ServerHttpsPort: "443",
		Db: DB{
			Server:   "192.168.1.63",
			Port:     1433,
			User:     "sts",
			Password: "sts",
			Database: "go_test",
		},
	}

	b, err := json.Marshal(&cj)
	if err != nil {
		logger.Error("Ошибка создания файла конфигурации", err)
	}

	return b
}

func createFileConfig(fp string) (*os.File, error) {
	f, err := os.Create(fp)
	if err != nil {
		return nil, err
	}
	f.Write(defaultJsonBytes())
	return f, nil
}

func openFileConfig() *os.File {
	var fP = "../config.json"
	f, err := os.OpenFile(fP, os.O_RDWR, 0755)
	if err != nil {
		if f, err = createFileConfig(fP); err != nil {
			logger.Error("2", "Ошибка создания файла конфигурации", err)
		}
	}
	return f
}
