package db

import (
	"database/sql"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"github.com/qzzznan/notification/model"
	log "github.com/sirupsen/logrus"
	"time"
)

func InsertUser(appleID, email, name, uuid string) error {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto(UserTable)
	ib.Cols("apple_id", "email", "name", "uuid", "create_at")
	ib.Values(appleID, email, name, uuid, time.Now())
	ib.SQL("ON CONFLICT (apple_id) DO NOTHING")

	str, args := ib.Build()

	log.Infoln("InsertUser:", str, args)

	_, err := db.Exec(str, args...)
	return err
}

func ExistUser(appleID string) (string, error) {
	s := sqlbuilder.PostgreSQL.NewSelectBuilder()
	s.Select("uuid").From(UserTable)
	s.Where(s.Equal("apple_id", appleID))
	str, args := s.Build()

	log.Infoln("ExistUser:", str, args)

	var uuid sql.NullString
	err := db.QueryRowx(str, args...).Scan(&uuid)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return uuid.String, nil
}

func GetUser(uuid, appleID string) (*model.User, error) {
	s := sqlbuilder.PostgreSQL.NewSelectBuilder()
	s.Select("id", "apple_id", "email", "name", "uuid", "create_at")
	s.From(UserTable)
	if uuid != "" {
		s.Where(s.Equal("uuid", uuid))
	} else if appleID != "" {
		s.Where(s.Equal("apple_id", appleID))
	} else {
		return nil, fmt.Errorf("uuid or apple_id is required")
	}
	query, args := s.Build()

	log.Infoln("GetUser:", query, args)

	user := &model.User{}
	err := db.Get(user, query, args...)
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
