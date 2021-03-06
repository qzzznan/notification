package model

type RegInfo struct {
	Token    string `form:"token" validate:"uuid4"`
	Name     string `form:"name" validate:"required"`
	DeviceID string `form:"device_id" validate:"required"`
	IsClip   int    `form:"is_clip"`
}

type APNsMessage struct {
	DeviceToken string
	Category    string
	Title       string
	Body        string
	Sound       string
	Data        map[string]interface{}
}
