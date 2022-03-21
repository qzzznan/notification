package db

func GetUserID(uuid string) (int64, error) {
	u, err := GetUser(uuid, "")
	if err != nil {
		return 0, err
	}
	return u.ID, nil
}
