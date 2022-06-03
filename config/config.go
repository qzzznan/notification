package config

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Http `yaml:"http"`
		Log  `yaml:"log"`
		PG   `yaml:"pg"`
		RDB  `yaml:"rdb"`
	}

	Http struct {
		Addr string `yaml:"addr" env:"HTTP_ADDR" env-default:"127.0.0.1:8006"`
	}

	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	}

	PG struct {
		User     string `yaml:"user" env:"PG_USER" env-default:"postgres"`
		Password string `yaml:"password" env:"PG_PWD" env-required:"true"`
		URL      string `yaml:"url" env:"PG_URL"`
		Port     string `yaml:"port" env:"PG_PORT" env-default:"5432"`
	}

	RDB struct {
		Addr     string `yaml:"addr" env:"RDB_ADDR" env-default:"127.0.0.1:6379"`
		Password string `yaml:"password" env:"RDB_PWD"`
	}
)

func NewConfig() (*Config, error) {
	c := &Config{}
	err := cleanenv.ReadConfig("config/config.yml", c)
	if err != nil {
		return nil, err
	}
	if c.Log.Level == "debug" {
		spew.Dump(c)
	}
	return c, err
}
