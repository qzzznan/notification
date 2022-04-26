package db

import (
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/huandu/go-assert"
	"github.com/qzzznan/notification/model"
	"github.com/qzzznan/notification/util"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	err := InitPostgresDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//clearDB()
	gofakeit.Seed(0)
	m.Run()
}

var inUsers []*model.User
var outUsers []*model.User

func TestInsertUser(t *testing.T) {
	n := 5
	inUsers = make([]*model.User, 0, n)
	for i := 0; i < n; i++ {
		appleID := gofakeit.LetterN(32)
		email := gofakeit.Email()
		name := gofakeit.Name()
		uuid := util.GenUID()

		inUsers = append(inUsers, &model.User{
			AppleID: appleID,
			Email:   email,
			Name:    name,
			UUID:    uuid,
		})
		err := InsertUser(appleID, email, name, uuid)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestExistUser(t *testing.T) {
	for _, v := range inUsers {
		exist, err := ExistUser(v.AppleID)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, exist, v.UUID)
	}
}

func TestGetUser(t *testing.T) {
	for _, v := range inUsers {
		user, err := GetUser(v.UUID, "")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v.AppleID, user.AppleID)
		outUsers = append(outUsers, user)
	}
}

func TestInsertDevice(t *testing.T) {
	for _, v := range outUsers {
		userID := v.ID
		deviceID := gofakeit.LetterN(32)
		typ := gofakeit.LetterN(3)
		name := gofakeit.PetName()
		device := &model.Device{
			UserID:   userID,
			DeviceID: deviceID,
			Type:     typ,
			Name:     name,
			IsClip:   0,
		}
		err := InsertDevice(device)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetDevice(t *testing.T) {
	for _, v := range outUsers {
		_, err := GetDevice(v.ID)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
	}
}

var devices []*model.Device

func TestGetAllDevice(t *testing.T) {
	for _, v := range outUsers {
		d, err := GetAllDevice(v.ID)
		if err != nil {
			t.Fatal(err)
		}
		devices = append(devices, d...)
	}
}

func TestUpdateDeviceName(t *testing.T) {
	for _, v := range devices {
		name := gofakeit.PetName()
		err := UpdateDeviceName(v.ID, name)
		if err != nil {
			t.Fatal(err)
		}
		device, err := GetDevice(v.ID)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, name, device.Name)
	}
}

func TestInsertPushKey(t *testing.T) {
	for _, v := range outUsers {
		name := gofakeit.PetName()
		key := gofakeit.LetterN(32)
		k := &model.PushKey{
			UserID: v.ID,
			Name:   name,
			Key:    key,
		}
		err := InsertPushKey(k)
		if err != nil {
			t.Fatal(err)
		}
	}
}

var keys []*model.PushKey

func TestGetAllPushKey(t *testing.T) {
	keys = make([]*model.PushKey, 0, len(outUsers))
	for _, v := range outUsers {
		ks, err := GetAllPushKey(v.ID)
		if err != nil {
			t.Fatal(err)
		}
		keys = append(keys, ks...)
	}
}

func TestGetPushKey(t *testing.T) {
	for _, v := range keys {
		key, err := GetPushKey(v.ID, v.Name)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v.UserID, key.UserID)
		assert.Equal(t, v.Name, key.Name)
		assert.Equal(t, v.Key, key.Key)
	}
}

func TestUpdatePushKey(t *testing.T) {
	for _, v := range keys {
		name := gofakeit.PetName()
		k := gofakeit.LetterN(32)
		err := UpdatePushKey(v.ID, name, k)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestAddMessage(t *testing.T) {
	for _, v := range outUsers {
		msg := &model.Message{
			UserID:  v.ID,
			Text:    gofakeit.Adjective(),
			Type:    gofakeit.StreetName(),
			Note:    gofakeit.PetName(),
			PushKey: gofakeit.LetterN(32),
			URL:     gofakeit.URL(),
			SendAt:  time.Now(),
		}
		err := AddMessage(msg)
		if err != nil {
			t.Fatal(err)
		}
	}
}
