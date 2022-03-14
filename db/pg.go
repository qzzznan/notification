package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
)

type PostgresDB struct {
	ctx context.Context
	*sql.DB
	sqlbuilder.Flavor
}

func (s *PostgresDB) Init(ctx context.Context) (err error) {
	s.Flavor = sqlbuilder.PostgreSQL
	s.ctx = ctx

	dsn := fmt.Sprintf("postgres://postgres:%s@%s/postgres?sslmode=disable",
		os.Getenv("PG_PWD"),
		os.Getenv("PG_URL"))

	log.Infoln("Connecting to Postgres:", dsn)

	s.DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return
	}

	return s.checkTables()
}

func (s *PostgresDB) checkTables() error {
	err := s.createTokenTable()
	if err != nil {
		return fmt.Errorf("create token table %v", err)
	}
	return nil
}

func (s *PostgresDB) createTokenTable() error {
	cre := s.NewCreateTableBuilder().CreateTable(TokenTable)
	cre.IfNotExists()
	cre.Define("user_token", "varchar(128)", "NOT NULL", "PRIMARY KEY")
	cre.Define("login_token", "varchar(512)", "NOT NULL")
	str, args := cre.Build()
	log.Infoln("create table:", str, args)
	_, err := s.ExecContext(s.ctx, str, args...)
	return err
}

func (s *PostgresDB) SaveLoginToken(userToken, loginToken string) error {
	ins := s.NewInsertBuilder().InsertInto(TokenTable)
	ins.Cols("user_token", "login_token")
	ins.Values(userToken, loginToken)
	str, args := ins.Build()

	log.Infoln("insert sql:", str, args)

	_, err := s.ExecContext(s.ctx, str, args...)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresDB) GetLoginToken(userToken string) (string, error) {
	sel := s.NewSelectBuilder().Select("login_token")
	sel.From(TokenTable)
	sel.Where(sel.Equal("user_token", userToken))
	str, args := sel.Build()

	log.Infoln("select sql:", str, args)

	row := s.QueryRowContext(s.ctx, str, args...)

	var loginToken string
	return loginToken, row.Scan(&loginToken)
}

func (s *PostgresDB) GetUserToken(loginToken string) (string, error) {
	sel := s.NewSelectBuilder().Select("user_token")
	sel.From(TokenTable)
	sel.Where(sel.Equal("login_token", loginToken))
	str, args := sel.Build()

	log.Infoln("select sql:", str, args)

	row := s.QueryRowContext(s.ctx, str, args...)

	var userToken string
	return userToken, row.Scan(&userToken)
}

func (s *PostgresDB) RemoveToken(userToken, loginToken string) error {
	del := s.NewDeleteBuilder().DeleteFrom(TokenTable)
	del.Where(del.Or(del.Equal("user_token", userToken), del.Equal("login_token", loginToken)))
	str, args := del.Build()

	log.Infoln("delete row:", str, args)

	_, err := s.ExecContext(s.ctx, str, args...)
	return err
}
