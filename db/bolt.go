package db

import (
	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
)

var bdb *bolt.DB

func InitBoltDB() (err error) {
	m := viper.GetStringMapString("boltdb")
	bdb, err = bolt.Open(m["path"], 0600, nil)
	return
}
