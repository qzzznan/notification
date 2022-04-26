package db

import "github.com/qzzznan/notification/model"

func GetUserID(token string) (int64, error) {
	u, err := GetUser(token, "")
	if err != nil {
		return 0, err
	}
	return u.ID, nil
}

func GetPushKeyInfo(pushKey string) (*model.PushKey, error) {
	return GetPushKey(0, "", pushKey)
}
