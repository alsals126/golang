package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"restapi/app" // restapi/app => 다른 패키지 import하기
	"sort"
	"strings"
)

// 처음화면
func menu(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>MAIN</h1>"+
		"<input type=\"button\" value=\"회원가입\" onClick=\"location.href='/join'\" /><br>"+
		"<input type=\"button\" value=\"회원조회\" onClick=\"location.href='/list'\" /><br>"+
		"<input type=\"button\" value=\"회원삭제\" onClick=\"location.href='/delete'\" /><br>"+
		"<input type=\"button\" value=\"회원수정\" onClick=\"location.href='/updateid'\" /><br>")
}

// func errView(err error, w http.ResponseWriter, errMSG string) bool {
// 	if err != nil {
// 		fmt.Fprintln(w, errMSG)
// 		return true
// 	}
// }

// 회원가입
func join(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "<h1>JOIN</h1>"+
		"<form action=\"/joinserver\" method=\"POST\">"+
		"ID <input type=\"text\" name=\"id\"/><br>"+
		"PW <input type=\"password\" name=\"pw\"/><br>"+
		"NAME <input type=\"text\" name=\"name\"/><br>"+
		"<input type=\"submit\" value=\"LOGIN\"> "+
		"<input type=\"button\" value=\"MAIN\" onClick=\"location.href='/'\">"+
		"</form>")
}
func joinServer(w http.ResponseWriter, r *http.Request) {
	errMSG := "<h1>회원가입에 실패하였습니다.</h1>" + "<input type=\"button\" value=\"MAIN\" onClick=\"location.href='/'\">"

	err := r.ParseForm()
	if err != nil || r.Form.Get("id") == "" {
		fmt.Fprintln(w, errMSG)
		return
	}

	params := url.Values{}
	params.Add("id", r.Form.Get("id"))
	params.Add("pw", r.Form.Get("pw"))
	params.Add("name", r.Form.Get("name"))

	// Post로 요청
	resp, err := http.PostForm("http://127.0.0.1:8080/restapi/join", params)
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	defer resp.Body.Close()

	// 응답 확인
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	var p app.Post
	err = json.Unmarshal([]byte(respBody), &p)
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	if p.State == "Success" {
		fmt.Fprintf(w, "<b>%s</b>님, 환영합니다.\t"+
			"<input type=\"button\" value=\"MAIN\" onClick=\"location.href='/'\">",
			p.Message["userid"])
	} else {
		fmt.Fprintln(w, errMSG)
	}
}

// 회원조회
func list(w http.ResponseWriter, _ *http.Request) {
	errMSG := "<h1>회원조회에 실패하였습니다.<h1> <input type=\"button\" value=\"MAIN\" onClick=\"location.href='/'\">"

	// GET 호출
	resp, err := http.Get("http://127.0.0.1:8080/restapi/list")
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	defer resp.Body.Close()

	// 응답 확인
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	// json unmarshal
	var p app.ListPost
	err = json.Unmarshal([]byte(respBody), &p)
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	if p.State == "Success" {
		// Map의 특정 key로 정렬
		sortKeys := make([]string, 0, len(p.Message))
		for k := range p.Message {
			sortKeys = append(sortKeys, p.Message[k].Userid)
		}
		sort.Strings(sortKeys)

		s := "<h1>LIST</h1><table border=\"1\"><th style=\"min-width: 150px\">USER ID</th><th style=\"min-width: 150px\">USER NAME</th>"
		for _, val := range p.Message {
			s += "<tr><td>" + val.Userid + "</td><td>" + val.Username + "</td></tr>"
		}
		s += "</table><br>" + "<input type=\"button\" value=\"MAIN\" onClick=\"location.href='/'\">"

		fmt.Fprintln(w, s)
	} else {
		fmt.Fprintln(w, errMSG)
	}
}

// 회원삭제
func delete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>DELETE</h1>"+
		"<form action=\"/deleteServer\" method=\"POST\">"+
		"ID <input type=\"text\" name=\"id\"/><br>"+
		"<input type=\"submit\" value=\"DELETE\"> "+
		"<input type=\"button\" value=\"MAIN\" onClick=\"location.href='/'\">"+
		"</form>")
}
func deleteServer(w http.ResponseWriter, r *http.Request) {
	errMSG := "<h1>회원삭제에 실패하였습니다.<h1>" + "<input type=\"button\" value=\"MAIN\" onClick=\"location.href='/'\">"

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	// DELETE 요청
	req, err := http.NewRequest("DELETE", "http://127.0.0.1:8080/restapi/delete?id="+r.Form.Get("id"), nil)
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}
	defer resp.Body.Close()

	// 응답 확인
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	var p app.Post
	err = json.Unmarshal([]byte(respBody), &p)
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	if p.State == "Success" && len(p.Message) > 0 {
		fmt.Fprintf(w, "<b>%s</b>님, 반가웠습니다..\t"+
			"<input type=\"button\" value=\"MAIN\" onClick=\"location.href='/'\">",
			p.Message["userid"])
	} else {
		fmt.Fprintln(w, errMSG)
	}
}

// 회원수정
func updateID(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>UPDATE</h1>"+
		"<form action=\"/update\" method=\"POST\">"+
		"ID <input type=\"text\" name=\"id\"/><br>"+
		"<input type=\"submit\" value=\"CHECK\"> "+
		"<input type=\"button\" value=\"MAIN\" onClick=\"location.href='/'\">"+
		"</form>")
}
func update(w http.ResponseWriter, r *http.Request) {
	errMSG := "<h1>회원수정에 실패하였습니다.\t</h1>" + "<input type=\"button\" value=\"MAIN\" onClick=\"location.href='/'\">"

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	// GET 호출
	resp, err := http.Get("http://127.0.0.1:8080/restapi/updatelist?id=" + r.Form.Get("id"))
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	defer resp.Body.Close()

	// 응답 확인
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	// 결과 출력
	var p app.Post
	err = json.Unmarshal([]byte(respBody), &p)
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	if p.State == "Success" && len(p.Message) > 0 {
		val := strings.Split(p.Message[r.Form.Get("id")], ",")

		fmt.Fprintf(w, "<h1>UPDATE</h1>"+
			"<form action=\"/updateServer\" method=\"POST\">"+
			"ID <input type=\"text\" name=\"id\" value=\""+val[0]+"\" readonly/><br>"+
			"PW <input type=\"text\" name=\"pw\" value=\""+val[1]+"\" /><br>"+
			"NAME <input type=\"text\" name=\"name\" value=\""+val[2]+"\" /><br>"+
			"<input type=\"submit\" value=\"UPDATE\"> "+
			"<input type=\"button\" value=\"MAIN\" onClick=\"location.href='/'\">"+
			"</form>")
	} else {
		fmt.Fprintln(w, errMSG)
	}
}
func updateServer(w http.ResponseWriter, r *http.Request) {
	errMSG := "<h1>회원수정에 실패하였습니다.\t</h1>" + "<input type=\"button\" value=\"MAIN\" onClick=\"location.href='/'\">"

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	params := url.Values{}
	params.Add("id", r.Form.Get("id"))
	params.Add("newPw", r.Form.Get("pw"))
	params.Add("newName", r.Form.Get("name"))

	// Post로 요청
	resp, err := http.PostForm("http://127.0.0.1:8080/restapi/update", params)
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	defer resp.Body.Close()

	// 응답 확인
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	var p app.Post
	err = json.Unmarshal([]byte(respBody), &p)
	if err != nil {
		fmt.Fprintln(w, errMSG)
		return
	}

	if p.State == "Success" {
		fmt.Fprintf(w, "<b>%s</b>님, 수정완료했습니다.\t"+
			"<input type=\"button\" value=\"MAIN\" onClick=\"location.href='/'\">",
			p.Message["userid"])
	} else {
		fmt.Fprintln(w, errMSG)
	}
}

func main() {
	log.Println("main.go started") //테스트코드

	http.HandleFunc("/", menu)
	http.HandleFunc("/join", join)
	http.HandleFunc("/joinserver", joinServer)
	http.HandleFunc("/list", list)
	http.HandleFunc("/delete", delete)
	http.HandleFunc("/deleteServer", deleteServer)
	http.HandleFunc("/updateid", updateID)
	http.HandleFunc("/update", update)
	http.HandleFunc("/updateServer", updateServer)
	app.ApiServer()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
