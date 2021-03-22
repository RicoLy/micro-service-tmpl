package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"micro-service-tmpl/utils/log"
	"micro-service-tmpl/utils/viper"
	"time"
)

type Configs struct {
	JwtSecret     string // Jwt秘钥
	SuperAdminId  int64  // 超级管理员Id
	SuperAdmin    string // 超级管理员账号
	SuperAdminPwd string // 超级管理员密码
}

var jwtConf Configs

func init() {
	if err := viper.ViperConf.UnmarshalKey("Jwt", &jwtConf); err != nil {
		log.GetLogger().Fatal("数据库获取配置文件失败" + err.Error())
	}
}

//Claims ...
type Claims struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Authority int    `json:"authority"`
	jwt.StandardClaims
}

//GenerateToken 签发用户Token
func GenerateToken(username, password string, authority int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)

	claims := Claims{
		username,
		password,
		authority,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "cmall",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(jwtConf.JwtSecret))

	return token, err
}

// ParseToken 验证用户token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConf.JwtSecret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

func ReNewJWT(claims *Claims) (string, error) {
	//更新JWT
	return GenerateToken(claims.Username, claims.Password, claims.Authority)
}

//EmailClaims ...
type EmailClaims struct {
	UserID        uint   `json:"user_id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	OperationType uint   `json:"operation_type"`
	jwt.StandardClaims
}

//GenerateEmailToken 签发邮箱验证Token
func GenerateEmailToken(userID, Operation uint, email, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(15 * time.Minute)

	claims := EmailClaims{
		userID,
		email,
		password,
		Operation,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "cmall",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(jwtConf.JwtSecret))

	return token, err
}

// ParseEmailToken 验证邮箱验证token
func ParseEmailToken(token string) (*EmailClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &EmailClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConf.JwtSecret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*EmailClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
