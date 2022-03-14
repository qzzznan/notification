package db

import (
	"github.com/qzzznan/notification/model"
	"github.com/qzzznan/notification/util"
)

func GetUserInfo() *model.User {
	return nil
}

func GetToken(UUID string) string {
	return util.GenToken(UUID)
}
