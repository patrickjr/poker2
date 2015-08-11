package users

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/gorilla/sessions"
	"time"
	//"fmt"
)

type User struct {
	Name      	string 
	Email     	string
	Online 		bool
	Remember 	bool
	IP 			string
}
/*
func (user *User) SetName		(name 	string) 	{user.name = name}
func (user *User) SetEmail		(email 	string) 	{user.email = email}	
func (user *User) SetOnline		(val 	bool) 		{user.online = val}
func (user *User) SetRemember	(val 	bool) 		{user.remember = val}
func (user *User) SetIP			(ip 	string) 	{user.ip = ip}

func (user *User) Name() 		string 	{return user.name}
func (user *User) Email() 		string 	{return user.email}
func (user *User) Online() 		bool   	{return user.online}
func (user *User) Remember() 	bool 	{return user.remember}
func (user *User) IP() 			string 	{return user.ip}
*/
func LoginWithCookie(db *sql.DB, name, c, ip string) (error){
	err := validateUserByCookie(db, name, c, ip)
	return err
}

func (user *User) Logout(db *sql.DB){
	logout(user, db)
}

func Default() (*User) {return &User{Online: false, Remember: false}}

//New generates a new user, pending the name and email are unique/valid
func New(db *sql.DB, name, email, passwordDigest string) error {

	err := validateName(db, name)
	if err != nil {
		return err
	}
	err = validateEmail(db, email)
	if err != nil {
		return err
	}
	passwordDigest, err = digestPassword(passwordDigest)
	if err != nil {
		return err
	}
	createdAt := time.Now().UTC()
	activationDigest, err := activationCode(passwordDigest)
	if err != nil {
		return errors.New("Internal Error AC_45")
	}
	stmt, err := db.Prepare("INSERT INTO users (name, email, password_digest, created_at, updated_at, activation_digest, activated) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(name, email, passwordDigest, createdAt, createdAt, activationDigest, 0)
	if err != nil {
		return err
	}
	return nil
}

//Login validates the user
func Login(db *sql.DB, username, password, rememberMe, myIP string) (*User, error) {
	
	user, err := validateUser(db, username, password)
	if err != nil {
		return nil, err
	}
	if rememberMe == "true"{
		user.Remember = true
		//user.SetRemember(true)
	}
	user.IP = myIP
	//user.SetIP(myIP)
	return user, nil
	
}




