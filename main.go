package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func userLogin(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()
	mobile := request.PostForm.Get("mobile")
	passwd := request.PostForm.Get("passwd")
	loginok := false

	if (mobile == "186" && passwd == "123") {

		loginok = true
	}

	if loginok {
		data := make(map[string]interface{})
		data["id"] = 1
		data["token"] = "test"
		Resp(writer,0, data, "")
	}else {
		Resp(writer, -1, nil,"password wrong")
	}

}

type H struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Resp(writer http.ResponseWriter,code int, data interface{}, msg string)  {
	//set header
	writer.Header().Set("Content-Type","application/json")
	writer.WriteHeader(http.StatusOK)

	//define a struct
	h := H{
		Code:code,
		Msg:msg,
		Data:data,
	}

	//transform the h to string
	ret,err := json.Marshal(h)

	if err!=nil {
		log.Println(err.Error())
	}
	
	writer.Write(ret)
}

func main() {

	//bind the func and request
	http.HandleFunc("/user/login",userLogin)

	//start the web server
	http.ListenAndServe(":8080",nil)

}
