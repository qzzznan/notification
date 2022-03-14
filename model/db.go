package model

import "time"

type User struct {
	ID       int64
	UUID     string
	IDToken  string
	CreateAt time.Time
	UpdateAt time.Time
}

type Device struct {
	ID       int64
	OwnUUID  string
	DeviceID string
	Type     string
	IsClip   string
	Name     string
	CreateAt time.Time
	UpdateAt time.Time
}

type PushKey struct {
	ID       int64
	OwnUUID  string
	Name     string
	Key      string
	CreateAt time.Time
	UpdateAt time.Time
}

type Message struct {
	ID      int64
	OwnUUID string
	Text    string
	Type    string
	Desc    string
	PushKey string
	URL     string
	SendAt  time.Time
}
