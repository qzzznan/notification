package db

import (
	"fmt"
	"github.com/qzzznan/notification/log"
	"github.com/qzzznan/notification/model"
)

var barkKeyCache map[string]*model.BarkDevice
var barkTokenCache map[string]*model.BarkDevice

func init() {
	barkKeyCache = make(map[string]*model.BarkDevice)
	barkTokenCache = make(map[string]*model.BarkDevice)
}

func GetUserIDStr(token string) (string, error) {
	u, err := GetUser(token, "")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", u.ID), nil
}

func GetPushKeyInfo(pushKey string) (*model.PushKey, error) {
	return GetPushKey(0, "", pushKey)
}

func InsertUser(appleID, email, name, uuid string) error {
	return insertUser(appleID, email, name, uuid)
}

func ExistUser(appleID string) (string, error) {
	return getUserIDByAppleID(appleID)
}

func GetUser(uid, appleID string) (*model.User, error) {
	return getUser(uid, appleID)
}

func InsertDevice(device *model.Device) error {
	return insertDevice(device)
}

func GetDevice(id int64) (*model.Device, error) {
	return getDevice(id)
}

func GetAllDevice(userID string) ([]*model.Device, error) {
	return getAllDevice(userID)
}

func UpdateDeviceName(id int64, newName string) error {
	return updateDeviceName(id, newName)
}

func InsertPushKey(key *model.PushKey) error {
	return insertPushKey(key)
}

func GetPushKey(id int64, name, pushKey string) (*model.PushKey, error) {
	return getPushKey(id, name, pushKey)
}

func GetAllPushKey(userID string) ([]*model.PushKey, error) {
	return getAllPushKey(userID)
}

func UpdatePushKey(keyID int64, newName, newKey string) error {
	return updatePushKey(keyID, newName, newKey)
}

func AddMessage(msg *model.Message) error {
	return addMessage(msg)
}

func GetMessages(userID string, offset, count int) ([]*model.Message, error) {
	return getMessages(userID, offset, count)
}

func RemoveDevice(id string) error {
	return removeDevice(id)
}

func RemoveKey(kid string) error {
	return removeKey(kid)
}

func RemoveMessage(id string) error {
	return removeMessage(id)
}

// ************************* Bark Begin *************************

func InsertBarkDevice(key, token string) error {
	m := &model.BarkDevice{
		DeviceKey:   key,
		DeviceToken: token,
	}
	barkTokenCache[token] = m
	barkKeyCache[key] = m
	return insertBarkDevice(m)
}

func GetBarkDevice(key, token string) (*model.BarkDevice, error) {
	if m, ok := barkKeyCache[key]; ok {
		return m, nil
	}
	if m, ok := barkTokenCache[token]; ok {
		return m, nil
	}
	m, err := getBarkDevice(key, token)
	if err == nil {
		barkKeyCache[m.DeviceKey] = m
		barkTokenCache[m.DeviceToken] = m
	}
	return m, err
}

func GetBarkToken(deviceKey string) string {
	d, err := GetBarkDevice(deviceKey, "")
	if err != nil {
		log.Errorf("GetBarkToken: %v", err)
		return ""
	}
	return d.DeviceToken
}

// ************************* Bark End ***************************
