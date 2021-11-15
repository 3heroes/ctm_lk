package config

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"ctm_lk/pkg/logger"

	"github.com/caarlos0/env/v6"
)

var Cfg *Config
var once sync.Once

type Config struct {
	servAddr     string
	dbConnString string
	appDir       string
}

func (c Config) ServAddr() string {
	return c.servAddr
}

func (c Config) DBConnString() string {
	return c.dbConnString
}

func (c Config) ProgramPath() string {
	return c.appDir
}

type EnvOptions struct {
	ServAddr     string `env:"RUN_ADDRESS"`
	DBConnString string `env:"DATABASE_URI"`
}

//checkEnv for get options from env to default application options.
func (c *Config) checkEnv() {

	e := new(EnvOptions)
	err := env.Parse(e)
	if err != nil {
		logger.Info("Ошибка чтения конфигурации из переменного окружения", err)
	}
	if len(e.ServAddr) != 0 {
		c.servAddr = e.ServAddr
	}
	if len(e.DBConnString) != 0 {
		c.dbConnString = e.DBConnString
	}

}

func (c *Config) setDefault() {
	c.servAddr = "localhost:8082"
	c.dbConnString = "user=postgres password=112233 dbname=ctm sslmode=disable"
}

//setFlags for get options from console to default application options.
func (c *Config) setFlags() {
	flag.StringVar(&c.servAddr, "a", c.servAddr, "a server address string")
	flag.StringVar(&c.dbConnString, "d", c.dbConnString, "a db connection string")
	flag.Parse()
}

func createConfig() {
	Cfg = new(Config)

	appDir, err := os.Getwd()
	if err != nil {
		logger.Error(err)
	}

	Cfg.setDefault()
	Cfg.checkEnv()
	Cfg.setFlags()
	Cfg.appDir = appDir
	fmt.Println(Cfg.DBConnString())
	logger.Info("Создан config")
}

// NewDefOptions return obj like Options interfase.
func NewConfig() {
	once.Do(createConfig)
}
