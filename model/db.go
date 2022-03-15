package model

import "time"

type User struct {
	ID       int64     `db:"id"`
	UUID     string    `db:"uuid"`
	IDToken  string    `db:"id_token"`
	CreateAt time.Time `db:"create_at"`
	UpdateAt time.Time `db:"update_at"`
}

type Device struct {
	ID       int64     `db:"id"`
	OwnUUID  string    `db:"own_uuid"`
	DeviceID string    `db:"device_id"`
	Type     string    `db:"type"`
	IsClip   string    `db:"is_clip"`
	Name     string    `db:"name"`
	CreateAt time.Time `db:"create_at"`
	UpdateAt time.Time `db:"update_at"`
}

type PushKey struct {
	ID       int64     `db:"id"`
	OwnUUID  string    `db:"own_uuid"`
	Name     string    `db:"name"`
	Key      string    `db:"key"`
	CreateAt time.Time `db:"create_at"`
	UpdateAt time.Time `db:"update_at"`
}

type Message struct {
	ID      int64     `db:"id"`
	OwnUUID string    `db:"own_uuid"`
	Text    string    `db:"text"`
	Type    string    `db:"type"`
	Note    string    `db:"note"`
	PushKey string    `db:"push_key"`
	URL     string    `db:"url"`
	SendAt  time.Time `db:"send_at"`
}
