package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// ErrWrongPassword : a
var ErrWrongPassword error = errors.New("wrong password")

// ErrWrongLoginOrPassword : a
var ErrWrongLoginOrPassword error = errors.New("wrong login or password")

// ErrAlreadyUsedLogin : a
var ErrAlreadyUsedLogin error = errors.New("already used login")

// Userstorage : a
type Userstorage interface {
	Valid() error
	Signup() error
}

// Signup : c
func Signup(username, pass string) error {
	err := Valid(username, pass)
	fmt.Println(err)
	if err != nil && err != ErrWrongPassword {
		salt := strconv.Itoa(rand.Intn(5000))
		data := username + " " + salt + " " + GetMD5hash(salt, pass) + "\n"
		file, err := os.OpenFile("users.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			return err
		}

		fmt.Println(data)
		_, err = file.WriteString(data)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		return nil

	}
	return ErrAlreadyUsedLogin
}

// Valid : b
func Valid(login string, pass string) error {
	file, err := ioutil.ReadFile("users.txt")
	if err != nil {
		return err
	}

	fileLines := strings.Split(string(file), "\n")
	for i := 0; i < len(fileLines); i++ {
		if fileLines[i] != "" {
			fileLine := strings.Split(string(fileLines[i]), " ")
			if fileLine[0] == login {
				if GetMD5hash(fileLine[1], pass) == fileLine[2] {
					return nil
				}
				return ErrWrongPassword
			}
		}
	}
	return ErrWrongLoginOrPassword
}

// GetMD5hash : f
func GetMD5hash(salt, pass string) string {
	hash := md5.New()
	hash.Write([]byte(pass))
	hash.Write([]byte(salt))
	return hex.EncodeToString(hash.Sum(nil))
}
