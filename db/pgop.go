package db

import (
	"database/sql"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"github.com/qzzznan/notification/model"
	log "github.com/sirupsen/logrus"
	"time"
)

func InsertUser(uuid, token string) error {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto(UserTable)
	ib.Cols("uuid", "id_token", "create_at")
	ib.Values(uuid, token, time.Now())
	ib.SQL("ON CONFLICT (id_token) DO NOTHING")

	str, args := ib.Build()

	log.Infoln("InsertUser:", str, args)

	_, err := db.Exec(str, args...)
	return err
}

func ExistUser(uuid, token string) (string, error) {
	s := sqlbuilder.PostgreSQL.NewSelectBuilder()
	s.Select("uuid").From(UserTable)

	if uuid != "" {
		s.Where(s.Equal("uuid", uuid))
	} else if token != "" {
		s.Where(s.Equal("id_token", token))
	} else {
		return "", fmt.Errorf("uuid or token is empty")
	}

	str, args := s.Build()

	log.Infoln("ExistUser:", str, args)

	var UUID sql.NullString
	err := db.QueryRowx(str, args...).Scan(&UUID)
	if err != nil {
		return "", err
	}
	if UUID.Valid {
		return UUID.String, nil
	}
	return "", nil
}

func GetUser(uuid string) (*model.User, error) {
	s := sqlbuilder.PostgreSQL.NewSelectBuilder()
	s.Select("id", "uuid", "id_token", "create_at", "update_at")
	s.From(UserTable)
	s.Where(s.Equal("uuid", uuid))
	query, args := s.Build()

	log.Infoln("GetUser:", query, args)

	user := &model.User{}
	err := db.Select(user, query, args)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func InsertDevice(device *model.Device) error {
	return nil
}

func GetDevice(deviceID string) *model.Device {
	return nil
}

func GetAllDevice(ownUUID string) []*model.Device {
	return nil
}

func UpdateDeviceName(deviceID string, newName string) error {
	return nil
}

func InsertPushKey(key *model.PushKey) error {
	return nil
}

func GetPushKey(key string) *model.PushKey {
	return nil
}

func GetAllPushKey(ownUUID string) []*model.PushKey {
	return nil
}

func UpdatePushKeyName(key string, newName string) error {
	return nil
}

func UpdatePushKey(key string, newKey string) error {
	return nil
}

func AddMessage(msg *model.Message) error {
	return nil
}
