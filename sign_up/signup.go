package main

import (
	"fmt"
	"log"
	"net/http"

	lib "OkonmaV/lib"
	us "OkonmaV/userstorage"

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
	err := userstorage.SignUp(login, pass)
	if err != nil {
		//if err.Error() == us.ErrAlreadyUsedLogin.Error() || us.ErrShortPassword.Error() {
		fmt.Fprintf(w, err.Error())
		//}
		//fmt.Println(err)
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
	log.Fatal(http.ListenAndServe(":8082", nil))
}
