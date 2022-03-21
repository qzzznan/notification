package db

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/huandu/go-assert"
	"github.com/qzzznan/notification/model"
	"github.com/qzzznan/notification/util"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := InitPostgresDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	gofakeit.Seed(0)
	m.Run()
}

var users []*model.User

func TestInsertUser(t *testing.T) {
	n := 5
	users = make([]*model.User, 0, n)
	for i := 0; i < n; i++ {
		appleID := gofakeit.LetterN(32)
		email := gofakeit.Email()
		name := gofakeit.Name()
		uuid := util.GenUID()

		users = append(users, &model.User{
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
	for _, v := range users {
		exist, err := ExistUser(v.AppleID)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, exist, v.UUID)
	}
}

func TestGetUser(t *testing.T) {
	for _, v := range users {
		user, err := GetUser(v.UUID, "")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, v.AppleID, user.AppleID)
	}
}
