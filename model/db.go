package model

import (
	"time"
)

type User struct {
	ID       int64     `db:"id" json:"id"`
	AppleID  string    `db:"apple_id" json:"apple_id"`
	Email    string    `db:"email" json:"email"`
	Name     string    `db:"name" json:"name"`
	UUID     string    `db:"uuid" json:"-"`
	CreateAt time.Time `db:"create_at" json:"-"`
	UpdateAt time.Time `db:"update_at" json:"-"`
	Level    int       `json:"level"`
	WeChatID string `json:"wechat_id"`
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
	ID       int64     `db:"id" json:"id"`
	UserID   int64     `db:"user_id" json:"uid"`
	Name     string    `db:"name" json:"name	"`
	Key      string    `db:"key" json:"key"`
	CreateAt time.Time `db:"create_at" json:"create_at"`
	UpdateAt time.Time `db:"update_at" json:"update_at"`
}

type Message struct {
	ID      int64     `db:"id" json:"id"`
	UserID  int64     `db:"user_id"`
	Text    string    `db:"text" json:"text"`
	Type    string    `db:"type" json:"type"`
	Note    string    `db:"note" json:"-"`
	PushKey string    `db:"push_key" json:"-"`
	URL     string    `db:"url" json:"-"`
	SendAt  time.Time `db:"send_at" json:"create_at"`
}
