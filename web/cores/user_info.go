package cores

import "github.com/dgrijalva/jwt-go"

type UserInfo struct {
	jwt.StandardClaims
	ID       uint   `json:"id"`
	Username string `json:"username"`
}
