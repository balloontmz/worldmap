package boot

import (
	"gopkg.in/ini.v1"
	"gopkg.in/urfave/cli.v1"
)

var (
	Conf *Config
)

type (
	Config struct {
		Debug      bool
		DataFile   string
		ConfigFile string
		*Srv
	}
	Srv struct {
		Host string
		Port string
	}
)

func init() {
	Conf = Default()
}

// DefaultConfig get default config
func Default() *Config {

	return &Config{
		false,
		"./LocList.xml",
		"config.ini",
		&Srv{
			"127.0.0.1",
			"1323",
		},
	}

}

// LoadFromIni load config from ini override default config
func (config *Config) LoadFromIni() (err error) {
	return ini.MapTo(config, config.ConfigFile)
}

// Load load config from command line param
func (config *Config) Load(c *cli.Context) (err error) {

	if c.String("config") != "" {
		Conf.ConfigFile = c.String("config")
		if err = Conf.LoadFromIni(); err != nil {
			return
		}
	}

	if c.Bool("debug") {
		Conf.Debug = true
	}

	if c.String("port") != "" {
		Conf.Srv.Port = c.String("port")
	}

	if c.String("bind") != "" {
		Conf.Srv.Host = c.String("bind")
	}

	return
}
