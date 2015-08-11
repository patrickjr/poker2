package users

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	//"fmt"
)

func validateUserByCookie(db *sql.DB, name, c, ip string) (error){
	stmt, err := db.Prepare("SELECT ip FROM users WHERE name=?")
	if err!=nil{
		return err
	}
	row := stmt.QueryRow(name)
	x := ip
	row.Scan(&x)
	if ip == x{
		return nil
	}
	return nil
}


func validateName(db *sql.DB, name string) error {
	if len(name) == 0 {
		return errors.New("username cannot be empty")
	}
	if len(name) > 15 {
		return errors.New("username cannot be longer than 15 characters")
	}
	err := uniqueName(db, name)
	if err != nil {
		return err
	}
	return nil
}

func validateEmail(db *sql.DB, email string) error {
	if len(email) == 0 {
		return errors.New("email cannot be empty")
	}
	emailRegex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	re := regexp.MustCompile(emailRegex)
	value := re.MatchString(email)
	if !value {
		return errors.New("email is invalid")
	}
	err := uniqueEmail(db, email)
	if err != nil {
		return err
	}
	return nil
}

func digestPassword(password string) (string, error) {
	ps := []byte(password)
	digest, err := bcrypt.GenerateFromPassword(ps, 10)
	if err != nil {
		return "error", err
	}
	password = bytesToString(digest)
	return password, nil
}
func comparePassword(password string, hash string) error {
	ps := []byte(password)
	hs := []byte(hash)
	err := bcrypt.CompareHashAndPassword(hs, ps)
	if err != nil {
		return errors.New("invalid combination")
	}
	return nil
}
func uniqueName(db *sql.DB, name string) error {
	stmt, err := db.Prepare("SELECT name FROM users WHERE name=?")
	if err != nil {
		return err
	}
	row := stmt.QueryRow(name)
	err = row.Scan(&name)
	switch {
	case err == sql.ErrNoRows:
		return nil
	default:
		return errors.New("username/email is taken")
	}
}

func uniqueEmail(db *sql.DB, email string) error {
	stmt, err := db.Prepare("SELECT email FROM users WHERE email=?")
	if err != nil {
		return err
	}
	row := stmt.QueryRow(email)
	err = row.Scan(&email)
	switch {
	case err == sql.ErrNoRows:
		return nil
	default:
		return errors.New("username/email is taken")
	}
}

func validateUser(db *sql.DB, username, passwordDigest string) (*User, error) {
	stmt, err := db.Prepare("SELECT name, password_digest, email FROM users WHERE name=?")
	if err != nil {
		return nil, err
	}
	var email string
	row := stmt.QueryRow(username)
	password := passwordDigest
	err = row.Scan(&username, &passwordDigest, &email)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("unsuccesful login")
	default:
		err = nil
	}
	err = comparePassword(password, passwordDigest)
	if err != nil {
		return nil, err
	}
	return &User{Name: username, Online: true, Email: email}, nil
}

func activationCode(passwordDigest string) (string, error) {
	activationDigest, err := digestPassword(passwordDigest)
	if err != nil {
		return "error", errors.New("Internal Error AC_VALIDATE45")
	}
	return activationDigest, nil
}

func bytesToString(b []byte) string {
	s := string(b[:])
	return s
}
