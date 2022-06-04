package entity

type BarkDevice struct {
	DeviceKey   string `db:"device_key" redis:"device_key"`
	DeviceToken string `db:"device_token" redis:"device_token"`
}

type BarkHistory struct {
}

type RegInfo struct {
	Token    string `form:"token" validate:"uuid4"`
	Name     string `form:"name" validate:"required"`
	DeviceID string `form:"device_id" validate:"required"`
	IsClip   int    `form:"is_clip"`
	Type     string `form:"Type"`
}

type APNsMessage struct {
	DeviceToken string
	Category    string
	Title       string
	Body        string
	Sound       string
	Data        map[string]interface{}
}
