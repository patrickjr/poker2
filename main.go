package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/rp/pkg/context"
	_ "html/template"
	"database/sql"
	"github.com/gorilla/sessions"
	"github.com/rp/src/poker2/controllers"
	"github.com/rp/src/poker2/models/users"
	"net/http"
	"encoding/gob"
	"html/template"
)
type webApplication struct{
	db *sql.DB
	store *sessions.CookieStore
	templ *template.Template
}

func (webApp *webApplication) ServeHTTP(w http.ResponseWriter, req *http.Request){

}

func (webApp *webApplication) Data(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
  		
  		if _, present := context.GetOk(r, "db"); !present {
 			context.Set(r, "db", webApp.db)
		}
  		if _, present := context.GetOk(r, "store"); !present{
  			context.Set(r, "store", webApp.store)
  		}
  		if _, present := context.GetOk(r, "templ"); !present{
  			context.Set(r, "templ", webApp.templ)
  		}
		fn(w, r) // call original
  	}
}


//redirect calls the parameter function prior to redirecting to "path"
func redirect(path string, fn http.HandlerFunc) http.HandlerFunc{ 
  return func (w http.ResponseWriter, r *http.Request) {
    fn(w, r)
    http.Redirect(w, r, path, 302)
  }
}


func NewWebApplication() *webApplication{
	webApp := new(webApplication)
	webApp.init()
	return webApp
}

func (webApp *webApplication) init() {
	database, err := sql.Open("mysql", "root:pawpisfsu@/rielpoker")
	if err!= nil{
		panic(err)
	}
	gob.RegisterName("user", &users.User{})
	webApp.db = database
	webApp.store = sessions.NewCookieStore([]byte("something-very-secret"))
	webApp.templ = template.Must(template.ParseGlob("/home/leebc/go/src/github.com/rp/src/poker2/views/templates/*")).Funcs(
		template.FuncMap{
			"eq": func(x, b bool) bool{
					return x == b
			},
		})
}

func main() {
	webApp := NewWebApplication()
	http.HandleFunc("/"			, webApp.Data(controllers.HomePage))
	http.HandleFunc("/about"	, webApp.Data(controllers.About))
	http.HandleFunc("/login"	, webApp.Data(controllers.Login))
	http.HandleFunc("/register"	, webApp.Data(controllers.Register))
	http.HandleFunc("/logout"	, redirect("/", webApp.Data(controllers.Logout)))
	
	http.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir("views"))))
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.ListenAndServe(":8082", context.ClearHandler(http.DefaultServeMux))
}
