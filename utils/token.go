package utils

import (
	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

const (
	ExpireTime = 24 // 过期时间
	//Secret     = "infant_mom" // 加盐
)

var (
	Secret = []byte("data_lot") // 加盐
)

type MyClaims struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Mobile    string `json:"mobile"`
	Head      string `json:"head"`
	Meid      string `json:"meid"`
	PhoneDesc string `json:"phone_desc"`
	ServiceId int64    `json:"service_id"`
	jwt.StandardClaims
}

// 解析token
func ParseToken(token string) (*MyClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*MyClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

// 创建 Token 值
func GetToken(claims *MyClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(Secret)
	if err != nil {
		log.Error("生成 Token 失败: ", err)
		return "", err
	}
	return signedToken, nil
}
