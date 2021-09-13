package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// DB정보
const (
	DB_USER     = "jeongmin"
	DB_PASSWORD = "1234asdf#$"
	DB_NAME     = "jeongmindb"
)

var dbinfo = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
	DB_USER, DB_PASSWORD, DB_NAME)

//sql.Open() 실행 시, 실제 DB를 연결하는 것이 아님. 쿼리를 던지는 시점에서 이루어진다.
var db, _ = sql.Open("postgres", dbinfo) // _ 자리는 원래 err자리. err검사를 하고싶으면 변수명 넣기

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// 로그인
func login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>LOGIN</h1>"+
		"<form action=\"/loginResult\" method=\"POST\">"+
		"ID <input type=\"text\" name=\"id\"/><br>"+
		"PW <input type=\"password\" name=\"pw\"/><br>"+
		"<input type=\"submit\" value=\"LOGIN\">"+
		"</form>")
}

// 로그인 실행시
func loginResult(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM login")
	CheckError(err)
	defer rows.Close()

	r.ParseForm() // 이거 필수

	str := ""
	for rows.Next() {
		var id string
		var pw string

		err = rows.Scan(&id, &pw)
		CheckError(err)

		if r.Form.Get("id") == id {
			if r.Form.Get("pw") == pw {
				str = "로그인되었습니다."
				break
			} else {
				str = "비밀번호가 틀렸습니다."
			}
		} else {
			str = "아이디 정보가 없습니다."
		}
	}
	fmt.Fprintln(w, str)
}

func main() {
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world")
	})
	http.HandleFunc("/login", login)
	http.HandleFunc("/loginResult", loginResult)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
