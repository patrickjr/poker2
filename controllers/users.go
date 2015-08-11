package controllers

import (
	"net/http"
	"github.com/rp/pkg/context"
	"fmt"	
	//"encoding/gob"
	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
	"database/sql"
	"github.com/rp/src/poker2/models/users"
	"html/template"
)
var (
	store *sessions.CookieStore
	db *sql.DB
)

func getWebAppData(r *http.Request)(*sql.DB, *sessions.CookieStore, *template.Template){
	db := context.Get(r, "db").(*sql.DB)
	store := context.Get(r, "store").(*sessions.CookieStore)
  	t := context.Get(r, "templ").(*template.Template)
  	return db, store, t	
}

func loginWithCookie(db *sql.DB, name, c, ip string) (error){
	err := users.LoginWithCookie(db, name, c, ip)
	return err 
}

func checkSession(r *http.Request, db *sql.DB, s *sessions.CookieStore) (*users.User, bool){
	exists, _ := s.Get(r, "rp-session")
	if exists.IsNew {
		return nil, false
	}else{
		x := exists.Values["user"]
		if x !=nil{
			user := x.(*users.User)
			err  := loginWithCookie(db, user.Name, "", r.RemoteAddr)
			if err!=nil{
				return nil, false
			}
			return user, true
		}
		return nil, false
	}
}

func renderTemplate(w http.ResponseWriter, tName string, tType *template.Template, dataObj interface{}) {

	err := tType.ExecuteTemplate(w, tName+".html", dataObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func saveCookie(w http.ResponseWriter, s *sessions.Session, store *sessions.CookieStore) (string, error) {
	encoded, err := securecookie.EncodeMulti(s.Name(), s.Values, store.Codecs...)
	if err != nil {
		return "", err
	}
	cookie := sessions.NewCookie(s.Name(), encoded, s.Options)
	http.SetCookie(w, cookie)
	return cookie.Value, nil
}

func setSessionOpts(r *http.Request, user *users.User, s *sessions.Session){
	if user.Remember{
		s.Options.MaxAge =  86400 * 360
	}
	//s.Options.Domain = r.Host
	s.Options.Path = "/"
	s.Options.HttpOnly = true
} 

func login(w http.ResponseWriter, r *http.Request, db *sql.DB, store *sessions.CookieStore, t *template.Template){
	if r.Method == "POST" {
		r.ParseForm()
		username, password, remember := r.FormValue("user[name]"), r.FormValue("user[password]"), r.FormValue("user[remember_me]")
		User, err := users.Login(db, username, password, remember, r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		s, _ := store.Get(r, "rp-session")
		setSessionOpts(r, User, s)
		s.Values["user"] = User
		err = s.Save(r, w)
		if err!=nil{
			fmt.Println(err)
		}
		renderTemplate(w, "user_nav_info", t, User)
	}	
}

func Login(w http.ResponseWriter, r *http.Request) {
	db, s, t := getWebAppData(r)
	login(w, r, db, s, t)
}

func register(w http.ResponseWriter, r *http.Request, db *sql.DB){
	if r.Method == "POST" {
		r.ParseForm()
		name, email, password := r.FormValue("user[name]"), r.FormValue("email"), r.FormValue("user[password]")
		err := users.New(db, name, email, password)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}	
}


func about(w http.ResponseWriter, t *template.Template){
	renderTemplate(w, "about", t, nil)
}

func About(w http.ResponseWriter, r *http.Request){
	_, _, t := getWebAppData(r)
	about(w, t)		
}


func logout(w http.ResponseWriter, r *http.Request, db *sql.DB, store *sessions.CookieStore, t *template.Template){
	s, _ := store.Get(r, "rp-session")
	if s.IsNew{
		http.Redirect(w, r, "/", 200)
	}
	//s.Options.MaxAge = -1
	user := s.Values["user"].(*users.User)
	user.Logout(db)
	s.Values["user"] = users.Default
	http.SetCookie(w, &http.Cookie{Name: "rp-session", MaxAge: -1, Path: "/"})
	http.Redirect(w, r, "/", 200)
}

func Logout(w http.ResponseWriter, r *http.Request){
	db, s, t := getWebAppData(r)
	logout(w, r, db, s, t)
}

func Register(w http.ResponseWriter, r *http.Request) {
	db, _, _ := getWebAppData(r)
	register(w, r, db)
}

func HomePage(w http.ResponseWriter, r *http.Request){
	db, store, t := getWebAppData(r)
  	if user, ok := checkSession(r, db, store); ok{
  		renderTemplate(w, "index", t, user)
  		
  	}else{
  		renderTemplate(w, "index", t, users.Default())
  	}
}
