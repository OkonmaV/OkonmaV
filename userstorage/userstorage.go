package userstorage

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

// ErrWrongPassword : f
var ErrWrongPassword error = errors.New("wrong password")

// ErrWrongLoginOrPassword : u
var ErrWrongLoginOrPassword error = errors.New("wrong login or password")

// ErrAlreadyUsedLogin : c
var ErrAlreadyUsedLogin error = errors.New("already used login")

// ErrBadUid : k
var ErrBadUid error = errors.New("bad uid")

// ErrShortPassword : d
var ErrShortPassword error = errors.New("short password")

// ErrWrongLogin : o
var ErrBadLogin error = errors.New("bad login at id getting")

type UsValid interface {
	GetUid(login string) (string, error)
	Valid(login, pass string) error
}

type UsSignUp interface {
	SignUp(login, pass string)
}

type UsAuth interface {
	Check(uid, role string) (bool, error)
}
type UsReader interface {
	readUs(filename string) (*os.File, error)
}
type UsTxt struct {
	Filename      string // fileLine: "login salt md5 uid role"
	Location      string
	RolesFilename string // fileLine: "uid role"
	//RolesLocation string
}

func NewUsTxt(filename, location, rolesfilename string) *UsTxt {
	return &UsTxt{Filename: filename, Location: location, RolesFilename: rolesfilename}
}

func (storage *UsTxt) readUs(filename string) (*os.File, error) {
	return os.OpenFile((storage.Location + filename), os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)
}

func getUsData(i UsReader, filename string) ([]byte, error) {
	var data []byte
	file, err := i.readUs(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err = ioutil.ReadAll(file)
	return data, err
}

// getMD5hash : pass + salt to md5
func getMD5hash(salt, pass string) (string, error) {
	hash := md5.New()
	_, err := hash.Write([]byte(pass))
	if err != nil {
		return "", err
	}
	_, err = hash.Write([]byte(salt))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Valid : Authentification
func (storage *UsTxt) Valid(login string, pass string) error {
	file, err := getUsData(storage, storage.Filename)
	if err != nil {
		return err
	}
	fileLines := strings.Split(string(file), "\n")
	for i := 0; i < len(fileLines); i++ {
		if fileLines[i] != "" {
			fileLine := strings.Split(string(fileLines[i]), " ")
			if fileLine[0] == login { // check typed login
				hash, err := getMD5hash(fileLine[1], pass)

				if err != nil {
					return err
				}
				if hash == fileLine[2] { // check typed pass
					return nil
				}
				return ErrWrongPassword
			}
		}
	}
	return ErrWrongLoginOrPassword
}

// Check : check your priveleges
func (storage *UsTxt) Check(uid, role string) (bool, error) {
	rolesfile, err := getUsData(storage, storage.RolesFilename)
	if err != nil {
		return false, err
	}
	fileLines := strings.Split(string(rolesfile), "\n")
	for i := 0; i < len(fileLines); i++ {
		if fileLines[i] != "" {
			fileLine := strings.Split(string(fileLines[i]), " ")
			if fileLine[0] == uid {
				if err != nil {
					return false, err
				}
				if role == "banned" {
					return true, nil
				}
				return false, nil
			}
		}
	}
	return false, ErrBadUid
}

// SignUp : registration
func (storage *UsTxt) SignUp(login, pass string) error {
	takenUid := 0
	file, err := getUsData(storage, storage.Filename)
	if err != nil {
		return err
	}
	fileLines := strings.Split(string(file), "\n") // check new login for being takened
	for takenUid < len(fileLines) {
		if fileLines[takenUid] != "" {
			fileLine := strings.Split(string(fileLines[takenUid]), " ")
			if login == fileLine[0] {
				return ErrAlreadyUsedLogin
			}
		}
		takenUid++
	}
	if len(pass) > 0 { // check new pass and write userdata to storage

		salt := strconv.Itoa(rand.Intn(5000))
		hash, err := getMD5hash(salt, pass)
		if err != nil {
			return err
		}

		file, err := storage.readUs(storage.Filename) // write to userstorage
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.WriteString(login + " " + salt + " " + hash + " " + strconv.Itoa(takenUid) + "\n")
		if err != nil {
			return err
		}

		rolesfile, err := storage.readUs(storage.RolesFilename) // write to rolesstorage
		if err != nil {
			return err
		}
		defer rolesfile.Close()
		_, err = rolesfile.WriteString(strconv.Itoa(takenUid) + " " + "смерд\n")
		return err
	}
	return ErrShortPassword
}

// GetUid : fuck
func (storage *UsTxt) GetUid(login string) (string, error) {
	file, err := getUsData(storage, storage.Filename)
	if err != nil {
		return "", err
	}
	fileLines := strings.Split(string(file), "\n")
	for i := 0; i < len(fileLines); i++ {
		if fileLines[i] != "" {
			fileLine := strings.Split(string(fileLines[i]), " ")
			fmt.Println(fileLine[0] + " - " + login)
			if fileLine[0] == login { // check typed login
				return fileLine[3], nil
			}
		}
	}
	return "", ErrBadLogin
}
