package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var user = make(map[string]string)

// Claims : fuck
type Claims struct {
	Login string `json:"login"`
	IP    string `json:"ip"`
	jwt.StandardClaims
}

func handler(w http.ResponseWriter, r *http.Request) {
	jwtKey := []byte("so_secure" + r.UserAgent())
	//	ip := r.Header.Get("X-Forwarded-For")
	cookie, err := r.Cookie("auth")
	if err != nil {
		fmt.Println("no cookie")
		return
	}
	tokenString := cookie.Value
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		fmt.Println("-----error with parsing token string-----")
		fmt.Println(err)
		return
	}
	if !token.Valid {
		fmt.Println("-----invalid token-----")
		fmt.Println(err)
		return
	}
	_, err = w.Write([]byte(fmt.Sprintf("User: %s", claims.Login)))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	user["mark"] = "123"
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8083", nil))
}
