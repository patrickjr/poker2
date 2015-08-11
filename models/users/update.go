package users

import(
	"database/sql"
)

func logout(user *User, db *sql.DB) (error){
	user.Online = false
	user.Remember = false
	stmt, err := db.Prepare("UPDATE users SET remember=? WHERE name=?")
	if err!=nil{return err}
	_, err = stmt.Exec(0, user.Name)
	return err
}