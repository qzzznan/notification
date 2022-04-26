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
	ib.Cols("apple_id", "email", "name", "uuid", "created_at")
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
	s.Select("id", "apple_id", "email", "name", "uuid", "created_at")
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
	ins := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ins.InsertInto(DeviceTable)
	ins.Cols("user_id", "device_id", "type", "is_clip", "name", "created_at", "updated_at")
	ins.Values(device.UserID, device.DeviceID, device.Type, device.IsClip, device.Name, time.Now(), time.Now())
	str, args := ins.Build()

	log.Infoln("InsertDevice:", str, args)

	_, err := db.Exec(str, args...)
	if err != nil {
		return err
	}
	return nil
}

func GetDevice(id int64) (*model.Device, error) {
	s := sqlbuilder.PostgreSQL.NewSelectBuilder()
	s.Select("id", "user_id", "device_id", "type", "is_clip", "name", "created_at", "updated_at")
	s.From(DeviceTable)
	s.Where(s.Equal("id", id))
	str, args := s.Build()

	log.Infoln("GetDevice:", str, args)

	device := &model.Device{}
	err := db.Get(device, str, args...)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func GetAllDevice(userID int64) ([]*model.Device, error) {
	s := sqlbuilder.PostgreSQL.NewSelectBuilder()
	s.Select("id", "user_id", "device_id", "type", "is_clip", "name", "created_at", "updated_at")
	s.From(DeviceTable)
	s.Where(s.Equal("user_id", userID))
	str, args := s.Build()

	log.Infoln("GetAllDevice:", str, args)

	devices := make([]*model.Device, 0)
	err := db.Select(&devices, str, args...)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func UpdateDeviceName(id int64, newName string) error {
	update := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	update.Update(DeviceTable)
	update.Set(update.Assign("name", newName))
	update.Where(update.Equal("id", id))
	str, args := update.Build()

	log.Infoln("UpdateDeviceName:", str, args)

	_, err := db.Exec(str, args...)
	if err != nil {
		return err
	}
	return nil
}

func InsertPushKey(key *model.PushKey) error {
	ins := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ins.InsertInto(PushKeyTable)
	ins.Cols("user_id", "name", "key", "created_at", "updated_at")
	ins.Values(key.UserID, key.Name, key.Key, time.Now(), time.Now())
	str, args := ins.Build()

	log.Infoln("InsertPushKey:", str, args)

	_, err := db.Exec(str, args...)
	if err != nil {
		return err
	}
	return nil
}

func GetPushKey(id int64, name, pushKey string) (*model.PushKey, error) {
	s := sqlbuilder.PostgreSQL.NewSelectBuilder()
	s.Select("id", "user_id", "name", "key", "created_at", "updated_at")
	s.From(PushKeyTable)
	if id != 0 {
		s.Where(s.Equal("id", id))
	} else if name != "" {
		s.Where(s.Equal("name", name))
	} else if pushKey != "" {
		s.Where(s.Equal("key", pushKey))
	} else {
		return nil, fmt.Errorf("id or name is required")
	}
	str, args := s.Build()

	log.Infoln("GetPushKey:", str, args)

	key := &model.PushKey{}
	err := db.Get(key, str, args...)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func GetAllPushKey(userID int64) ([]*model.PushKey, error) {
	s := sqlbuilder.PostgreSQL.NewSelectBuilder()
	s.Select("id", "user_id", "name", "key", "created_at", "updated_at")
	s.From(PushKeyTable)
	s.Where(s.Equal("user_id", userID))
	str, args := s.Build()

	log.Infoln("GetAllPushKey:", str, args)

	list := make([]*model.PushKey, 0)
	return list, db.Select(&list, str, args...)
}

func UpdatePushKey(keyID int64, newName, newKey string) error {
	update := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	update.Update(PushKeyTable)
	if newName != "" {
		update.Set(update.Assign("name", newName))
	}
	if newKey != "" {
		update.Set(update.Assign("key", newKey))
	}
	update.Where(update.Equal("id", keyID))
	str, args := update.Build()

	log.Infoln("UpdatePushKeyName:", str, args)

	_, err := db.Exec(str, args...)
	return err
}

func AddMessage(msg *model.Message) error {
	ins := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ins.InsertInto(MessageTable)
	ins.Cols("user_id", "text", "type", "note", "push_key", "url", "send_at")
	ins.Values(msg.UserID, msg.Text, msg.Type, msg.Note, msg.PushKeyName, msg.URL, msg.SendAt)
	str, args := ins.Build()

	log.Infoln("AddMessage:", str, args)

	_, err := db.Exec(str, args...)
	return err
}

func clearDB() {
	var err error
	for _, v := range []string{
		DeviceTable, PushKeyTable, MessageTable, UserTable,
	} {
		_, err = db.Exec("TRUNCATE TABLE " + v)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
