package main

import (
	"fmt"
	"net/http"
)

// 회원가입
func join(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>JOIN</h1>"+
		"<form action=\"/loginResult\" method=\"POST\">"+
		"ID <input type=\"text\" name=\"id\"/><br>"+
		"NAME <input type=\"text\" name=\"name\"/><br>"+
		"PW <input type=\"password\" name=\"pw\"/><br>"+
		"한번 더 PW <input type=\"password\" name=\"temPw\"/><br>"+
		"<input type=\"submit\" value=\"LOGIN\">"+
		"</form>")
}

/*
func main() {
	app := infoApp{}
	defer app.dbConnection(
		"jeongmin",
		"1234asdf#$",
		"jeongmindb",
	)
}
*/
