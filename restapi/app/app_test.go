package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"restapi/app"
	"testing"
)

func TestJoin(t *testing.T) {
	log.Println("app_test.go \"Join\" Test.....")

	// 요청
	resp, err := http.PostForm("http://127.0.0.1:8080/restapi/join",
		url.Values{"id": {"userid"}, "name": {"username"}, "pw": {"userpw"}})
	if err != nil {
		t.Error("REQUEST ERROR")
	}

	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("RESPONSE ERROR")
	}

	var p app.Post
	err = json.Unmarshal([]byte(respBody), &p)
	if err != nil {
		t.Error("JSON UNMARSHAL ERROR")
	} else {
		if p.State == "Success" && p.Message == nil {
			t.Error("WRONG RESULT")
		}
	}
}
func TestList(t *testing.T) {
	log.Println("app_test.go \"List\" Test.....")

	// 요청
	resp, err := http.Get("http://127.0.0.1:8080/restapi/list")
	if err != nil {
		t.Error("REQUEST ERROR")
	}

	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("RESPONSE ERROR")
	}

	var p app.ListPost
	err = json.Unmarshal([]byte(respBody), &p)
	if err != nil {
		t.Error("JSON UNMARSHAL ERROR")
	} else {
		if p.State == "Success" && p.Message == nil {
			t.Error("WRONG RESULT")
		}
	}
}
func TestDelete(t *testing.T) {
	log.Println("app_test.go \"Delete\" Test.....")

	// 요청
	req, err := http.NewRequest("DELETE", "http://127.0.0.1:8080/restapi/delete?id="+"userid", nil)
	if err != nil {
		t.Error("REQUEST ERROR")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error("REQUEST ERROR")
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("RESPONSE ERROR")
	}

	var p app.Post
	err = json.Unmarshal([]byte(respBody), &p)
	if err != nil {
		fmt.Println(p)
		t.Error("JSON UNMARSHAL ERROR")
	} else {
		if p.State != "Success" && p.Message != nil {
			t.Error("WRONG RESULT")
		}
	}
}
func TestUpdate(t *testing.T) {
	log.Println("app_test.go \"Update\" Test.....")

	bodyData := url.Values{
		"id":      {"userid"},
		"newName": {"username1"},
		"newPw":   {"userpw1"},
	}

	// 요청
	req, err := http.NewRequest("PUT", "http://127.0.0.1:8080/restapi/update", bytes.NewBufferString(bodyData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Error("REQUEST ERROR")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error("REQUEST ERROR")
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("RESPONSE ERROR")
	}

	var p app.Post
	err = json.Unmarshal([]byte(respBody), &p)
	if err != nil {
		t.Error("JSON UNMARSHAL ERROR")
	} else {
		if p.State == "Success" && p.Message == nil {
			t.Error("WRONG RESULT")
		}
	}
}
