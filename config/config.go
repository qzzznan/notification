package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var (
	configPath string
	configName string
	configExt  string
)

func LoadConfig() {
	flag.StringVar(&configPath, "p", ".", "config file path")
	flag.StringVar(&configName, "n", "config", "config file name")
	flag.StringVar(&configExt, "e", "yaml", "config file extension")
	flag.Parse()

	viper.SetConfigName(configName)
	viper.SetConfigType(configExt)
	viper.AddConfigPath(configPath)
	viper.SetConfigPermissions(0666)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.SetDefault("http.addr", "localhost:8006")
	viper.SetDefault("http.debug", true)

	viper.SetDefault("boltdb.path", "bolt.db")

	viper.SetDefault("postgres.url", os.Getenv("PG_URL"))
	viper.SetDefault("postgres.password", os.Getenv("PG_PWD"))
	viper.SetDefault("postgres.clear_db", false)

	viper.SetDefault("redis.addr", "localhost:6379")

	viper.SetDefault("log.log_file_path", "/var/log/notification/notification.log")
	viper.SetDefault("log.pid_file_path", "/var/run/notification.pid")
	viper.SetDefault("log.level", "debug")

	viper.SetDefault("push_deer", "./static/c.p12")

	fmt.Println("config ************************************")
	for k, v := range viper.AllSettings() {
		fmt.Printf("config key:%s value:%v\n", k, v)
	}
	fmt.Println("config ************************************")
}
