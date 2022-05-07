package db

import (
	"github.com/qzzznan/notification/model"
	"testing"
)

func TestHash(t *testing.T) {
	err := InitRedis()
	if err != nil {
		t.Fatal(err)
	}

	pk := &model.PushKey{
		ID:     123,
		UserID: "AAA",
		Name:   "BBB",
		Key:    "CCC",
	}

	err = SetStruct("Hello", pk)
	if err != nil {
		t.Fatal(err)
	}

	ppk := &model.PushKey{}
	err = GetStruct("Hello", ppk)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ppk)
}
