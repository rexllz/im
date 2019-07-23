package main

import (
	"io"
	"net/http"
)

func main() {

	//bind the func and request
	http.HandleFunc("/user/login",
		func(writer http.ResponseWriter, request *http.Request) {
		io.WriteString(writer,"hello world")
	})

	//start the web server
	http.ListenAndServe(":8080",nil)


}
