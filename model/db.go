package model

import (
	"time"
)

type User struct {
	ID        int64     `db:"id" json:"id"`
	AppleID   string    `db:"apple_id" json:"apple_id"`
	Email     string    `db:"email" json:"email"`
	Name      string    `db:"name" json:"name"`
	UUID      string    `db:"uuid" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Level     int       `json:"level"`
	WeChatID  string    `json:"wechat_id"`
}

type Device struct {
	ID        int64     `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"uid"`
	DeviceID  string    `db:"device_id" json:"device_id"`
	Type      string    `db:"type" json:"type"`
	IsClip    int       `db:"is_clip" json:"is_clip"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

type PushKey struct {
	ID        int64     `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"uid"`
	Name      string    `db:"name" json:"name"`
	Key       string    `db:"key" json:"key"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

type Message struct {
	ID          int64     `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"uid"`
	Text        string    `db:"text" json:"text"`
	Type        string    `db:"type" json:"type"`
	Note        string    `db:"note" json:"desp"`
	PushKeyName string    `db:"push_key" json:"pushkey_name"`
	URL         string    `db:"url" json:"-"`
	SendAt      time.Time `db:"send_at" json:"created_at"`
}
