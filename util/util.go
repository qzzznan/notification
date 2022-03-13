package util

import "github.com/satori/go.uuid"

func GenUID() string {
	return uuid.NewV4().String()
}
