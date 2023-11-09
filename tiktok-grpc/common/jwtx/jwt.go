package jwtx

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type DouShengClaims struct {
	UserName string `json:"username"`
	UserID   int64  `json:"userid"`
	jwt.RegisteredClaims
}

func CreateUserToken(id int64, name string) string {

	key := []byte("keybyxs")

	claim := DouShengClaims{
		UserID:   id,
		UserName: name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                 // 签发者
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), //过期时间
			//NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), //最早使用时间
			//IssuedAt:  jwt.NewNumericDate(time.Now()),                  //签发时间
			//Subject:   "Tom",                                           // 签发对象
			//Audience:  jwt.ClaimStrings{"Android_APP", "IOS_APP"},      //签发受众
			//ID:        randStr(10),                                     // wt ID, 类似于盐值
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	s, err := t.SignedString(key)
	if err != nil {
		fmt.Println(err)
	}

	return s
}

func ParseToken(tokenString string) (DouShengClaims, bool) {
	if tokenString == "" {
		return DouShengClaims{}, false
	}

	key := []byte("keybyxs")

	token, _ := jwt.ParseWithClaims(tokenString, &DouShengClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	claims, ok := token.Claims.(*DouShengClaims)

	if !ok || !token.Valid {

		return DouShengClaims{}, false
	}

	return *claims, true
}
