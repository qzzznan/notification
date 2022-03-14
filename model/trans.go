package model

type RegInfo struct {
	Token    string `form:"token"`
	Name     string `form:"name"`
	DeviceID string `form:"device_id"`
	IsClip   int    `form:"is_clip"`
}
