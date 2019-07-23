package main

import (
	"net/http"
)

func main() {

	//bind the func and request
	http.HandleFunc("/user/login",
		func(writer http.ResponseWriter, request *http.Request) {

			request.ParseForm()
			mobile := request.PostForm.Get("mobile")
			passwd := request.PostForm.Get("passwd")
			loginok := false

			if (mobile == "186" && passwd == "123") {

				loginok = true
			}

			str := `{"code":0,"data":{"id":1,"token":"test"}}`

			if !loginok {

				str = `{"code":-1,"msg":"password wrong"}`
			}
			//set header
			writer.Header().Set("Content-Type","application/json")
			writer.WriteHeader(http.StatusOK)

			writer.Write([]byte(str))

			//io.WriteString(writer,"hello world")
	})

	//start the web server
	http.ListenAndServe(":8080",nil)

}
