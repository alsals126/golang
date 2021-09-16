package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

type infoApp struct {
	name string
	id   string
	pw   string
}
type Post struct {
	State   string
	Message map[string]string
}

var db *sql.DB

func dbConnection(db_user, db_password, db_name string) {
	var dbinfo = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		db_user, db_password, db_name)

	var err error
	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
}

var stateStr string = ""

// json으로 인코딩하는 미들웨어
func jsonEncoding(handler func(http.ResponseWriter, *http.Request) map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := handler(w, r)

		// json encoding
		w.Header().Set("Content-Type", "application/json")
		post := &Post{
			State:   stateStr,
			Message: msg,
		}
		json, _ := json.Marshal(post)
		w.Write(json)
	}
}

// 회원가입
func (i infoApp) joinHandler(w http.ResponseWriter, r *http.Request) map[string]string {
	err := r.ParseForm()
	if err != nil {
		stateStr = "FORM ERROR"
	}
	i = infoApp{
		r.Form.Get("name"),
		r.Form.Get("id"),
		r.Form.Get("pw"),
	}

	_, err = db.Exec("INSERT INTO restapi1(name, id, pw) VALUES($1, $2, $3)", i.name, i.id, i.pw)
	if err != nil {
		stateStr = "DB ERROR"
	} else {
		stateStr = "Success"
	}

	m := map[string]string{
		"userid": i.id,
	}

	return m
}

// 회원목록
func (i infoApp) listHandler(w http.ResponseWriter, r *http.Request) map[string]string {
	rows, err := db.Query("SELECT * FROM login")
	if err != nil {
		stateStr = "DB ERROR"
	}
	defer rows.Close()

	m := make(map[string]string)
	for rows.Next() {
		var id string
		var name string

		err = rows.Scan(&id, &name)
		m[id] = "userid: " + id + ", username: " + name
	}

	return m
}
func main() {
	dbConnection("jeongmin", "1234asdf#$", "jeongmindb")
	defer db.Close()

	pathPrefix := "/restapi"

	i := infoApp{}

	http.HandleFunc(pathPrefix, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API Server")
	})
	http.HandleFunc(pathPrefix+"/join", jsonEncoding(i.joinHandler))
	http.HandleFunc(pathPrefix+"/list", jsonEncoding(i.listHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
