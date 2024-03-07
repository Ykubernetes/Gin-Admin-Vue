package utils

import (
	"gitee.com/go-server/global"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	encryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		global.Log.Fatal("密码加密存在错误:", err)
		return "", err
	} else {
		return string(encryptPassword), nil
	}
}

func EqualsPassword(password, encryptPassword string) bool {
	// 第一个参数为加密后的密码，第二个参数为未加密的密码
	err := bcrypt.CompareHashAndPassword([]byte(encryptPassword), []byte(password))
	// 对比密码是否正确会返回一个异常，按照官方的说法是只要异常是 nil 就证明密码正确
	return err == nil
}
