package model

import (
	"time"
)

type User struct {
	ID       int64     `db:"id"`
	AppleID  string    `db:"apple_id"`
	Email    string    `db:"email"`
	Name     string    `db:"name"`
	UUID     string    `db:"uuid"`
	CreateAt time.Time `db:"create_at"`
	UpdateAt time.Time `db:"update_at"`
}

type Device struct {
	ID       int64     `db:"id" json:"id"`
	UserID   int64     `db:"user_id" json:"uid"`
	DeviceID string    `db:"device_id" json:"device_id"`
	Type     string    `db:"type" json:"type"`
	IsClip   string    `db:"is_clip" json:"is_clip"`
	Name     string    `db:"name" json:"name"`
	CreateAt time.Time `db:"create_at" json:"-"`
	UpdateAt time.Time `db:"update_at" json:"-"`
}

type PushKey struct {
	ID       int64     `db:"id"`
	UserID   int64     `db:"user_id"`
	Name     string    `db:"name"`
	Key      string    `db:"key"`
	CreateAt time.Time `db:"create_at"`
	UpdateAt time.Time `db:"update_at"`
}

type Message struct {
	ID      int64     `db:"id"`
	UserID  int64     `db:"user_id"`
	Text    string    `db:"text"`
	Type    string    `db:"type"`
	Note    string    `db:"note"`
	PushKey string    `db:"push_key"`
	URL     string    `db:"url"`
	SendAt  time.Time `db:"send_at"`
}
