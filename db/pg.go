package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/qzzznan/notification/log"
	"github.com/spf13/viper"
)

var db *sqlx.DB

func InitPostgresDB() error {
	m := viper.GetStringMapString("postgres")
	dsn := fmt.Sprintf("postgres://postgres:%s@%s/postgres?sslmode=disable",
		m["password"], m["url"])
	log.Infoln("Connecting to Postgres:", dsn)
	db1, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	db = sqlx.NewDb(db1, "postgres")

	ctx := context.Background()
	for _, v := range TableCreateMap {
		err = v(ctx, db1)
		if err != nil {
			return err
		}
	}

	if m["clear_db"] == "true" {
		err = clearDB()
		if err != nil {
			return err
		}
	}

	return nil
}

const (
	UserTable    = "t_user"
	DeviceTable  = "t_device"
	PushKeyTable = "t_push_key"
	MessageTable = "t_message"

	BarkDeviceTable = "t_bark_device"
)

var TableCreateMap = map[string]func(ctx context.Context, db *sql.DB) error{
	UserTable:       createUserTable,
	DeviceTable:     createDeviceTable,
	PushKeyTable:    createPushKeyTable,
	MessageTable:    createMessageTable,
	BarkDeviceTable: createBarkDeviceTable,
}

func createUserTable(ctx context.Context, db *sql.DB) error {
	ctb := sqlbuilder.PostgreSQL.NewCreateTableBuilder()
	ctb.CreateTable(UserTable).IfNotExists()
	ctb.Define("id", "SERIAL", "PRIMARY KEY")
	ctb.Define("apple_id", "VARCHAR(255)", "NOT NULL", "UNIQUE")
	ctb.Define("email", "VARCHAR(127)")
	ctb.Define("name", "VARCHAR(127)")
	ctb.Define("uuid", "UUID", "NOT NULL", "UNIQUE")
	ctb.Define("created_at", "TIMESTAMP", "NOT NULL")
	ctb.Define("updated_at", "TIMESTAMP")

	str, args := ctb.Build()

	log.Infoln("CreateUserTable:", str, args)

	_, err := db.ExecContext(ctx, str, args...)
	if err != nil {
		return fmt.Errorf("create table %s: %w", UserTable, err)
	}
	return nil
}

func createDeviceTable(ctx context.Context, db *sql.DB) error {
	ctb := sqlbuilder.PostgreSQL.NewCreateTableBuilder()
	ctb.CreateTable(DeviceTable).IfNotExists()
	ctb.Define("id", "SERIAL", "PRIMARY KEY")
	ctb.Define("user_id", "INT", "NOT NULL")
	ctb.Define("device_id", "VARCHAR(64)", "NOT NULL", "UNIQUE")
	ctb.Define("type", "VARCHAR(32)")
	ctb.Define("is_clip", "SMALLINT")
	ctb.Define("name", "VARCHAR(32)")
	ctb.Define("created_at", "TIMESTAMP", "NOT NULL")
	ctb.Define("updated_at", "TIMESTAMP")

	str, args := ctb.Build()

	log.Infoln("CreateDeviceTable:", str, args)

	_, err := db.ExecContext(ctx, str, args...)
	if err != nil {
		return fmt.Errorf("create table %s: %w", DeviceTable, err)
	}
	return nil
}

func createPushKeyTable(ctx context.Context, db *sql.DB) error {
	ctb := sqlbuilder.PostgreSQL.NewCreateTableBuilder()
	ctb.CreateTable(PushKeyTable).IfNotExists()
	ctb.Define("id", "SERIAL", "PRIMARY KEY")
	ctb.Define("user_id", "INT", "NOT NULL")
	ctb.Define("name", "VARCHAR(32)")
	ctb.Define("key", "VARCHAR(64)", "NOT NULL", "UNIQUE")
	ctb.Define("created_at", "TIMESTAMP", "NOT NULL")
	ctb.Define("updated_at", "TIMESTAMP")

	str, args := ctb.Build()

	log.Infoln("CreatePushKeyTable:", str, args)

	_, err := db.ExecContext(ctx, str, args...)
	if err != nil {
		return fmt.Errorf("create table %s: %w", PushKeyTable, err)
	}
	return nil
}

func createMessageTable(ctx context.Context, db *sql.DB) error {
	ctb := sqlbuilder.PostgreSQL.NewCreateTableBuilder()
	ctb.CreateTable(MessageTable).IfNotExists()
	ctb.Define("id", "SERIAL", "PRIMARY KEY")
	ctb.Define("user_id", "INT", "NOT NULL")
	ctb.Define("text", "TEXT")
	ctb.Define("type", "VARCHAR(32)")
	ctb.Define("note", "VARCHAR(64)")
	ctb.Define("push_key", "VARCHAR(64)", "NOT NULL")
	ctb.Define("url", "VARCHAR(255)")
	ctb.Define("send_at", "TIMESTAMP", "NOT NULL")

	str, args := ctb.Build()

	log.Infoln("CreateMessageTable:", str, args)

	_, err := db.ExecContext(ctx, str, args...)
	if err != nil {
		return fmt.Errorf("create table %s: %w", MessageTable, err)
	}
	return nil
}

func createBarkDeviceTable(ctx context.Context, db *sql.DB) error {
	ctb := sqlbuilder.PostgreSQL.NewCreateTableBuilder()
	ctb.CreateTable(BarkDeviceTable).IfNotExists()
	ctb.Define("device_key", "UUID", "PRIMARY KEY")
	ctb.Define("device_token", "VARCHAR(256)", "NOT NULL")

	str, args := ctb.Build()
	log.Infoln("CreateBarkDeviceTable:", str, args)

	_, err := db.ExecContext(ctx, str, args...)
	if err != nil {
		return fmt.Errorf("create table %s: %w", BarkDeviceTable, err)
	}
	return nil
}
