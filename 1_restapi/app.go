package main

import (
	"database/sql"
	"encoding/json"
	"errors"
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

// json으로 인코딩하는 미들웨어
func (p *Post) jsonEncoding(handler func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)

		// json encoding
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			p.State = err.Error()
		}
		fmt.Print(p)
		json, _ := json.Marshal(p)
		w.Write(json)
	}
}

// 회원가입
func (p *Post) joinHandler(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return errors.New("FORM ERROR")
	}
	i := infoApp{
		r.Form.Get("name"),
		r.Form.Get("id"),
		r.Form.Get("pw"),
	}

	_, err = db.Exec("INSERT INTO restapi1(name, id, pw) VALUES($1, $2, $3)", i.name, i.id, i.pw)
	if err != nil {
		return errors.New("DB ERROR")
	} else {
		p.State = "Success"
		p.Message = map[string]string{
			"userid": i.id,
		}
		return err
	}
}

// 회원목록
func (p *Post) listHandler(_ http.ResponseWriter, _ *http.Request) error {
	rows, err := db.Query("SELECT id, name FROM restapi1")
	if err != nil {
		return errors.New("DB ERROR1")
	}
	defer rows.Close()

	fmt.Println(rows)
	for rows.Next() {
		var id string
		var name string

		err = rows.Scan(&id, &name)
		p.Message[id] = "userid: " + id + ", username: " + name
	}

	p.State = "Success"
	return err
}

// 회원탈퇴
func (p *Post) deleteHandler(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return errors.New("FORM ERROR")
	}

	_, err = db.Exec("DELETE FROM restapi1 WHERE id=$1", r.FormValue("id"))
	if err != nil {
		return errors.New("DB ERROR")
	} else {
		p.State = "Success"
		p.Message = map[string]string{
			"userid": r.Form.Get("id"),
		}
		return err
	}
}

func main() {
	dbConnection("jeongmin", "1234asdf#$", "jeongmindb")
	defer db.Close()

	pathPrefix := "/restapi"
	p := Post{}

	http.HandleFunc(pathPrefix, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API Server")
	})
	http.HandleFunc(pathPrefix+"/join", p.jsonEncoding(p.joinHandler))
	http.HandleFunc(pathPrefix+"/list", p.jsonEncoding(p.listHandler))
	http.HandleFunc(pathPrefix+"/delete", p.jsonEncoding(p.deleteHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
