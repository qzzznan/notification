package db

import (
	"context"
	"testing"
)

func TestSQLiteDB(t *testing.T) {
	db := SQLiteDB{}
	err := db.Init(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	ut := "userToken1"
	lt := "loginToken1"

	err = db.SaveLoginToken(ut, lt)
	if err != nil {
		t.Fatal(err)
	}

	qt, err := db.GetUserToken(lt)
	if err != nil {
		t.Fatal(err)
	}
	if qt != ut {
		t.Fatal("result error, expect:", ut, "real:", qt)
	}

	qt, err = db.GetLoginToken(ut)
	if err != nil {
		t.Fatal(err)
	}
	if qt != lt {
		t.Fatal("result error, expect:", lt, "real:", qt)
	}

	err = db.RemoveToken(ut, "")
	if err != nil {
		t.Fatal(err)
	}

	qt, err = db.GetUserToken("213123")
	t.Log(err)
	t.Log("->", qt)
}
