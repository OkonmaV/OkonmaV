package main

import (
	lib "OkonmaV/lib"
	us "OkonmaV/userstorage"
	"fmt"
	"log"
	"net/http"

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

func handler(w http.ResponseWriter, r *http.Request) {
	userstorage := us.NewUsTxt("users.txt", "../userstorage/", "roles.txt")
	login := r.URL.Query().Get("login")
	pass := r.URL.Query().Get("pass")
	err := userstorage.Valid(login, pass)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	err = lib.CreateCookie(w, r, userstorage, login)
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, r.Header.Get("Referer"), 302) // redirect

}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8083", nil))
}
