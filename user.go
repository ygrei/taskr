package main

import (
	"database/sql"
	"fmt"
	"errors"
)

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
)

const (
	salt = "here be dragons"
)

var (
	errNoSuchUser = errors.New("no such user")
)

type User struct {
	ID             int64
	Username       string
	RealName       string
	GitHubUsername string
	GitHubToken    string
	Location       string
	Techs          []string `db:"-"` // gorp does not support slices :(, see https://github.com/coopernurse/gorp/issues/5
	TechsRaw       string
	Email          string
	PasswordHash   string
	Bio            string
	PastProjects   string

	// Disbursement options
	DisburseStyle string
	PayPal        string
	Bitcoin       string
}

func clear(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}

func hashPassword(password []byte) ([]byte, error) {
	defer clear(password)
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func getUserByHandle(handle string) (*User, error) {
	var userID *int64
	if strings.Contains(handle, "@") {
		err := dbMap.SelectOne(&userID, "SELECT ID FROM users WHERE email=?", handle)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("unexpected DB error in getUserByHandle(%v): %v", handle, err)
		}
	} else {
		err := dbMap.SelectOne(&userID, "SELECT ID FROM users WHERE username=?", handle)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("unexpected DB error in getUser(%v): %v", handle, err)
		}
	}
	if userID == nil {
		return nil, errNoSuchUser
	}
	return getUser(*userID)
}

func getUser(id int64) (*User, error) {
	var user *User
	err := dbMap.SelectOne(&user, "SELECT * FROM users WHERE id=?", id)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("unexpected DB error in getUser(%v): %v", id, err)
	}
	if user == nil {
		return nil, errNoSuchUser
	}
	user.Expand()
	return user, nil
}

func (u *User) Expand() {
	if u == nil || u.TechsRaw == "" {
		return
	}
	techsSlice := strings.Split(u.TechsRaw, ",")
	for i, _ := range techsSlice {
		techsSlice[i] = strings.TrimSpace(techsSlice[i])
	}
	u.Techs = techsSlice
	return
}
