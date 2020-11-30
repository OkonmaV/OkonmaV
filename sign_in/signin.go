package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims : fuck
type Claims struct {
	Login string `json:"login"`
	IP    string `json:"ip"`
	jwt.StandardClaims
}

var user = make(map[string]string)

func handler(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Forwarded-For")
	signlogin := r.URL.Query().Get("login")
	signpassword := r.URL.Query().Get("pass")

	originpassword, ok := user[signlogin]
	if ok {
		if originpassword == signpassword {
			expTime := time.Now().Add(10 * time.Minute)
			claims := &Claims{
				Login: signlogin,
				IP:    ip,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expTime.Unix(),
				},
			}
			jwtKey := []byte("so_secure" + r.UserAgent())
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				fmt.Println("-----Problem with token creating-----")
				fmt.Println(err)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:    "auth",
				Value:   tokenString,
				Expires: expTime,
			})
			http.Redirect(w, r, r.Header.Get("Referer"), 302) // redirect

		} else {
			fmt.Fprintf(w, "meh (wrong password)")
		}
	} else {
		fmt.Fprintf(w, "other meh (no user with this login)")
	}

}

func main() {
	user["mark"] = "123"
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
