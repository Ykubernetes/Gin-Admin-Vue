package utils

import (
	"context"
	"gitee.com/go-server/global"
	"github.com/mojocn/base64Captcha"
	"time"
)

type CaptchaStore struct {
}

var ctx = context.Background()

// 实现captcha存储的set方法
func (c CaptchaStore) Set(id, value string) error {
	key := "captcha:" + id
	err := global.RedisDB.Set(ctx, key, value, time.Minute*2).Err()
	return err
}

// 实现captcha存储的Get方法
func (c CaptchaStore) Get(id string, clear bool) string {
	key := "captcha:" + id
	val, err := global.RedisDB.Get(ctx, key).Result()
	if err != nil {
		global.Log.Warnln("Redis获取验证码失败，%s", err)
		return ""
	}
	if clear {
		err := global.RedisDB.Del(ctx, key).Err()
		if err != nil {
			global.Log.Fatalf("Redis删除验证码失败，%s", err)
			return ""
		}
	}
	return val
}

func (c CaptchaStore) Verify(id, answer string, clear bool) bool {
	val := c.Get(id, clear)
	return val == answer
}

var store = CaptchaStore{}

// 生成图形化数字验证码配置
func digitConfig() *base64Captcha.DriverDigit {
	digitType := &base64Captcha.DriverDigit{
		Height:   50,
		Width:    120,
		Length:   6,
		MaxSkew:  0.5,
		DotCount: 78,
	}
	return digitType
}

func CreateCode() (string, string, error) {
	var driver base64Captcha.Driver
	// 创建验证码并传入创建的类型的配置，以及存储的对象
	driver = digitConfig()
	c := base64Captcha.NewCaptcha(driver, store)
	id, s, _, err := c.Generate() // id 验证码id // bse64s 图片base64编码 // err 错误
	return id, s, err

}

// 校验验证码
func VerifyCaptcha(id, verifyValue string) bool {
	return store.Verify(id, verifyValue, true)
}

// 获取验证码答案
func GetCodeAnser(CodeId string) string {
	return store.Get(CodeId, false)
}
