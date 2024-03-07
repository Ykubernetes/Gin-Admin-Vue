package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserInfo struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`

	jwt.RegisteredClaims // v5版本新加的方法
}

// 生成JWT
func GenerateJWT(username, uuid, secretKey string) (string, error) {
	claims := UserInfo{
		username,
		uuid,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 过期时间24小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                     // 生效时间
		},
	}
	// 使用HS256签名算法
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString([]byte(secretKey))

	return s, err
}

// 解析JWT
func ParseJwt(tokenstring, secretKey string) (*UserInfo, error) {
	t, err := jwt.ParseWithClaims(tokenstring, &UserInfo{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if claims, ok := t.Claims.(*UserInfo); ok && t.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
