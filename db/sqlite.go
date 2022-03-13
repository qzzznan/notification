package db

import (
	"context"
	"database/sql"
	"github.com/huandu/go-sqlbuilder"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type SQLiteDB struct {
	ctx context.Context
	*sql.DB
}

func (s *SQLiteDB) Init(ctx context.Context) (err error) {
	s.ctx = ctx
	s.DB, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		return
	}
	cre := sqlbuilder.CreateTable(TokenTable)
	cre.IfNotExists()
	cre.Define("user_token", "varchar(128)", "NOT NULL", "PRIMARY KEY")
	cre.Define("login_token", "varchar(512)", "NOT NULL")
	str, args := cre.Build()

	log.Infoln("create table:", str, args)

	_, err = s.ExecContext(ctx, str, args...)
	return
}

func (s *SQLiteDB) SaveLoginToken(userToken, loginToken string) error {
	ins := sqlbuilder.InsertInto(TokenTable)
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

func (s *SQLiteDB) GetLoginToken(userToken string) (string, error) {
	sel := sqlbuilder.Select("login_token")
	sel.From(TokenTable)
	sel.Where(sel.Equal("user_token", userToken))
	str, args := sel.Build()

	log.Infoln("select sql:", str, args)

	row := s.QueryRowContext(s.ctx, str, args...)

	var loginToken string
	return loginToken, row.Scan(&loginToken)
}

func (s *SQLiteDB) GetUserToken(loginToken string) (string, error) {
	sel := sqlbuilder.Select("user_token")
	sel.From(TokenTable)
	sel.Where(sel.Equal("login_token", loginToken))
	str, args := sel.Build()

	log.Infoln("select sql:", str, args)

	row := s.QueryRowContext(s.ctx, str, args...)

	var userToken string
	return userToken, row.Scan(&userToken)
}

func (s *SQLiteDB) RemoveToken(userToken, loginToken string) error {
	del := sqlbuilder.DeleteFrom(TokenTable)
	del.Where(del.Or(del.Equal("user_token", userToken), del.Equal("login_token", loginToken)))
	str, args := del.Build()

	log.Infoln("delete row:", str, args)

	_, err := s.ExecContext(s.ctx, str, args...)
	return err
}
