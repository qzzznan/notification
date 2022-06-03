package app

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func Migrate(dsn, opt string) {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	if opt == "up" {
		err = m.Up()
	} else if opt == "down" {
		err = m.Down()
	} else {
		err = errors.New("invalid option")
	}
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalln(err)
	}
	se, de := m.Close()
	if se != nil {
		log.Fatalln(se)
	}
	if de != nil {
		log.Fatalln(de)
	}
}
