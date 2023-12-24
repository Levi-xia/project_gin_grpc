package base

import (
	"time"

	"com.levi/project-common/config"
	"github.com/dgrijalva/jwt-go"
)

type jwtService struct {
}

var JwtService = new(jwtService)

type CustomClaims struct {
    jwt.StandardClaims
}

const (
    TokenType = "bearer"
    AppGuardName = "app"
)

type TokenOutPut struct {
    AccessToken string `json:"access_token"`
    ExpiresIn int `json:"expires_in"`
    TokenType string `json:"token_type"`
}

func (jwtService *jwtService) CreateToken(GuardName string, userId string) (tokenData TokenOutPut, err error) {
    token := jwt.NewWithClaims(
        jwt.SigningMethodHS256,
        CustomClaims{
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: time.Now().Unix() + config.GlobalConf.Jwt.JwtTtl,
                Id:        userId,
                Issuer:    GuardName, // 用于在中间件中区分不同客户端颁发的 token，避免 token 跨端使用
                NotBefore: time.Now().Unix() - 1000,
            },
        },
    )
    tokenStr, err := token.SignedString([]byte(config.GlobalConf.Jwt.Secret))
    tokenData = TokenOutPut{
        tokenStr,
        int(config.GlobalConf.Jwt.JwtTtl),
        TokenType,
    }
    return
}