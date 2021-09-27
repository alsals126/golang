// main 패키지와 main함수를 사용하기 위해 따로 만듬
package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	_ "github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

// User struct
type InfoApp struct {
	name string
	id   string
	pw   string
}

// response struct
type Post struct {
	State   string
	Message map[string]string
}
type ListPost struct {
	State   string
	Message []List
}
type List struct {
	Userid   string
	Username string
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
func jsonEncoding(handler func(http.ResponseWriter, *http.Request) (Post, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := handler(w, r)

		// json encoding
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			p.State = err.Error()
		}
		json, _ := json.Marshal(p)
		w.Write(json)
	}
}
func jsonEncodingList(handler func(http.ResponseWriter, *http.Request) (ListPost, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := handler(w, r)

		// json encoding
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			p.State = err.Error()
		}
		json, _ := json.Marshal(p)
		w.Write(json)
	}
}

// 회원가입
func joinHandler(w http.ResponseWriter, r *http.Request) (Post, error) {
	p := Post{}

	err := r.ParseForm()
	if err != nil {
		return p, errors.New("FORM ERROR")
	}
	i := InfoApp{
		r.Form.Get("name"),
		r.Form.Get("id"),
		r.Form.Get("pw"),
	}

	_, err = db.Exec("INSERT INTO restapi1(name, id, pw) VALUES($1, $2, $3)", i.name, i.id, i.pw)
	if err != nil {
		return p, errors.New("DB ERROR")
	} else {
		p = Post{
			State: "Success",
			Message: map[string]string{
				"userid": i.id,
			},
		}
		return p, err
	}
}

// 회원목록
func listHandler(w http.ResponseWriter, r *http.Request) (ListPost, error) {
	p := ListPost{}

	rows, err := db.Query("SELECT id, name FROM restapi1 order by id")
	if err != nil {
		return p, errors.New("DB ERROR1")
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string

		err := rows.Scan(&id, &name)
		if err != nil {
			return p, errors.New("DB ERROR2")
		}
		p.Message = append(p.Message, List{
			Userid:   id,
			Username: name,
		})
	}
	p.State = "Success"

	return p, err
}

// 유저 정보가 있는지
func ifuser(id string) error {
	result, err := db.Exec("SELECT * FROM restapi1 WHERE id=$1", id)
	if err != nil {
		return errors.New("DB ERROR")
	} else {
		nRow, err := result.RowsAffected()
		if err != nil {
			return errors.New("DB ERROR")
		}

		if nRow == 0 {
			return errors.New("NO USER")
		}
		return err
	}
}

// 회원탈퇴
func deleteHandler(w http.ResponseWriter, r *http.Request) (Post, error) {
	p := Post{}

	err := ifuser(r.URL.Query()["id"][0])
	if err != nil {
		return p, err
	}
	_, err = db.Exec("DELETE FROM restapi1 WHERE id=$1", r.URL.Query()["id"][0])
	if err != nil {
		return p, errors.New("DB ERROR")
	} else {
		p = Post{
			State: "Success",
			Message: map[string]string{
				"userid": r.URL.Query()["id"][0],
			},
		}
		return p, err
	}
}

// 회원수정
func updateHandler(w http.ResponseWriter, r *http.Request) (Post, error) {
	p := Post{}

	err := r.ParseForm()
	if err != nil {
		return p, errors.New("FORM ERROR")
	}
	i := InfoApp{
		r.Form.Get("newName"),
		r.Form.Get("id"),
		r.Form.Get("newPw"),
	}

	_, err = db.Exec("UPDATE restapi1 SET name=$1, pw=$2 WHERE id=$3", i.name, i.pw, i.id)
	if err != nil {
		return p, errors.New("DB ERROR")
	} else {
		p = Post{
			State: "Success",
			Message: map[string]string{
				"userid": i.id,
			},
		}
		return p, err
	}
}

// 회원모든정보
func updateList(w http.ResponseWriter, r *http.Request) (Post, error) {
	p := Post{}
	msg := make(map[string]string)

	err := ifuser(r.URL.Query()["id"][0])
	if err != nil {
		return p, err
	}

	rows, err := db.Query("SELECT * FROM restapi1 where id=$1", r.URL.Query()["id"][0])
	if err != nil {
		return p, errors.New("DB ERROR1")
	}
	defer rows.Close()

	for rows.Next() {
		var id, pw, name string

		err := rows.Scan(&name, &pw, &id)
		if err != nil {
			return p, errors.New("DB ERROR2")
		}
		msg[id] = id + "," + pw + "," + name
	}

	p = Post{
		State:   "Success",
		Message: msg,
	}
	return p, err
}

// func main() {
// 	log.Println("app.go/ApiServer started") //테스트코드

// 	dbConnection("jeongmin", "1234asdf#$", "jeongmindb")
// 	defer db.Close()

// 	pathPrefix := "/restapi"
// 	http.HandleFunc(pathPrefix, func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintln(w, "API Server")
// 	})
// 	http.HandleFunc(pathPrefix+"/join", jsonEncoding(joinHandler))
// 	http.HandleFunc(pathPrefix+"/list", jsonEncodingList(listHandler))
// 	http.HandleFunc(pathPrefix+"/delete", jsonEncoding(deleteHandler))
// 	http.HandleFunc(pathPrefix+"/update", jsonEncoding(updateHandler))
// 	http.HandleFunc(pathPrefix+"/updatelist", jsonEncoding(updateList))

// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
