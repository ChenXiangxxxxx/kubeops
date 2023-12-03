package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/wonderivan/logger"
)

var JwtToken jwtToken

type jwtToken struct{}

// 定义token反序列化后的内容，一般会放用户信息
type CustomClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// 加解密因子
const (
	SECRET = "ops" //在前端的login.vue以及路由的index.js中均用到。
)

func (*jwtToken) ParseToken(tokenString string) (claims *CustomClaims, err error) {
	//使用这个方法结合加密因子解析token，这个token是前端传过来的。获得一个Token类型的对象。
	//主要是吧tokenString序列化成CustomClaims结构体。
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	if err != nil {
		logger.Error("parse token failed ", err)
		//处理token解析后的各种错误
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("TokenMalformed")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("TokenExpired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("TokenNotValidYet")
			} else {
				return nil, errors.New("TokenInvalid")
			}
		}
	}
	//将Token对象中的Claims断言成CustomClaims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("解析Token无效")
}
