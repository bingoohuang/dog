package main

import (
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

type Conf struct {
	Programs map[string]Program
}

type Program struct {
	Cmd          string
	MaxMem       string
	Timeout      string
	PoolSize     int
	ExpectPrefix string
}

func LoadConf(configFile string) (Conf, error) {
	var config Conf
	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		logrus.Errorf("DecodeFile error %v", err)
		return config, err
	}

	return config, nil
}

func MustLoadConf(configFile string) Conf {
	config, err := LoadConf(configFile)
	if err != nil {
		logrus.Panic(err)
	}

	logrus.Debugf("config: %+v\n", config)
	return config
}
