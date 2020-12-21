package lib

import (
	us "OkonmaV/userstorage"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims : fuck
type Claims struct {
	Login     string
	IP        string
	UserAgent string
	Uid       string
	jwt.StandardClaims
}

// CreateCookie : <
func CreateCookie(w http.ResponseWriter, r *http.Request, userstorage *us.UsTxt, login string) error {
	ip := r.Header.Get("X-Real-IP")
	expTime := time.Now().Add(10 * time.Minute)
	uid, err := userstorage.GetUid(login)
	if err != nil {
		fmt.Println(err)
		return err
	}
	claims := &Claims{
		Login:     login,
		IP:        ip,
		UserAgent: r.UserAgent(),
		Uid:       uid,
	}
	jwtKey := []byte("so_secure")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "koki",
		Value:   tokenString,
		Expires: expTime,
	})
	return nil
}
