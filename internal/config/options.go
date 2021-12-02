package config

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"ctm_lk/pkg/logger"
)

var Cfg *Config
var once sync.Once

type Config struct {
	servHost      string
	servHttpPort  string
	servHttpsPort string
	dbConnString  string
	appDir        string
	dbServ        msSQLparams
}

type msSQLparams struct {
	server   string
	port     int
	user     string
	password string
	database string
}

func (c Config) ServAddrHttps() string {
	return net.JoinHostPort(c.servHost, c.servHttpsPort)
}

func (c Config) ServAddrHttp() string {
	return net.JoinHostPort(c.servHost, c.servHttpPort)
}

func (c Config) DBConnString() string {
	return c.dbConnString
}

func (c Config) ProgramPath() string {
	return c.appDir
}

func (c *Config) setDefault() {
	c.dbServ = msSQLparams{
		server:   "192.168.1.63",
		port:     1433,
		user:     "sts",
		password: "sts",
		database: "go_test",
	}
	c.servHost = "localhost"
	c.servHttpPort = "80"
	c.servHttpsPort = "443"
	// подключение к MS SQL Server
	c.dbConnString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		c.dbServ.server, c.dbServ.user, c.dbServ.password, c.dbServ.port, c.dbServ.database)

}

//setFlags for get options from console to default application options.
func (c *Config) getFlags() {
	flag.StringVar(&c.servHost, "a", c.servHost, "a server address string")
	flag.StringVar(&c.servHttpsPort, "p", c.servHttpsPort, "a https port string")
	flag.StringVar(&c.dbConnString, "d", c.dbConnString, "a db connection string")
	flag.Parse()
}

func (c *Config) checkServAddressValid() {
	if strings.Contains(c.servHost, ":") {
		if host, port, err := net.SplitHostPort(c.servHost); err == nil {
			c.servHost = host
			c.servHttpPort = port
			fmt.Println(host)
			fmt.Println(port)
		} else {
			logger.Info("Error", "Ошибка разбора адреса сервера. Не верный формат", c.servHost)
		}
	}
}

func createConfig() {
	Cfg = new(Config)

	appDir, err := os.Getwd()
	if err != nil {
		logger.Error(err)
	}
	f := openFileConfig()
	defer f.Close()
	fmt.Println(f)
	Cfg.setDefault()
	Cfg.getFlags()
	Cfg.appDir = appDir
	Cfg.checkServAddressValid()
	logger.Info("Создан config")
}

// NewDefOptions return obj like Options interfase.
func NewConfig() {
	once.Do(createConfig)
}
