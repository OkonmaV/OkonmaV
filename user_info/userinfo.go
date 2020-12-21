package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// Claims : fuck
type Claims struct {
	Login string
	IP    string
	Uid   string
	jwt.StandardClaims
}

func handler(w http.ResponseWriter, r *http.Request) {
	jwtKey := []byte("so_secure")
	//	ip := r.Header.Get("X-Forwarded-For")
	cookie, err := r.Cookie("koki")
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
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8084", nil))
}
