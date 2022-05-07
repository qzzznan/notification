package model

import (
	"time"
)

type User struct {
	ID        int64     `db:"id" json:"id" redis:"id"`
	AppleID   string    `db:"apple_id" json:"apple_id" redis:"apple_id"`
	Email     string    `db:"email" json:"email" redis:"email"`
	Name      string    `db:"name" json:"name" redis:"name"`
	UUID      string    `db:"uuid" json:"-" redis:"uuid"`
	CreatedAt time.Time `db:"created_at" json:"created_at" redis:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" redis:"-"`
	Level     int       `json:"level" redis:"level"`
	WeChatID  string    `json:"wechat_id" redis:"we_chat_id"`
}

type Device struct {
	ID        int64     `db:"id" json:"id" redis:"id"`
	UserID    string    `db:"user_id" json:"uid" redis:"user_id"`
	DeviceID  string    `db:"device_id" json:"device_id" redis:"device_id"`
	Type      string    `db:"type" json:"type" redis:"type"`
	IsClip    int       `db:"is_clip" json:"is_clip" redis:"is_clip"`
	Name      string    `db:"name" json:"name" redis:"name"`
	CreatedAt time.Time `db:"created_at" json:"-" redis:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-" redis:"-"`
}

type PushKey struct {
	ID        int64     `db:"id" json:"id" redis:"id"`
	UserID    string    `db:"user_id" json:"uid" redis:"user_id"`
	Name      string    `db:"name" json:"name" redis:"name"`
	Key       string    `db:"key" json:"key" redis:"key"`
	CreatedAt time.Time `db:"created_at" json:"created_at" redis:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-" redis:"-"`
}

type Message struct {
	ID          int64     `db:"id" json:"id" redis:"id"`
	UserID      string    `db:"user_id" json:"uid" redis:"user_id"`
	Text        string    `db:"text" json:"text" redis:"text"`
	Type        string    `db:"type" json:"type" redis:"type"`
	Note        string    `db:"note" json:"desp" redis:"note"`
	PushKeyName string    `db:"push_key" json:"pushkey_name" redis:"push_key_name"`
	URL         string    `db:"url" json:"-" redis:"url"`
	SendAt      time.Time `db:"send_at" json:"created_at" redis:"-"`
}

type BarkDevice struct {
	DeviceKey   string `db:"device_key" redis:"device_key"`
	DeviceToken string `db:"device_token" redis:"device_token"`
}

type BarkMessage struct {
}
