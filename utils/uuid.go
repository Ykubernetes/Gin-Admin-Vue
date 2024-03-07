package utils

import uuid "github.com/satori/go.uuid"

func GetUUID() string {
	return "userIdKey-" + uuid.NewV4().String()
}
