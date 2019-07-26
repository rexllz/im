package main

import (
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"im/ctrl"
	"log"
	"net/http"
)

func RegisterView()  {
	tpl,err := template.ParseGlob("view/**/*")
	if err!=nil {
		//quit and print the err
		log.Fatal(err.Error())
	}

	for _,v := range tpl.Templates(){
		tplname := v.Name()

		http.HandleFunc(tplname,
			func(writer http.ResponseWriter, request *http.Request) {
				tpl.ExecuteTemplate(writer, tplname, nil)
		})
	}
}



func main() {

	//bind the func and request
	http.HandleFunc("/user/login",ctrl.UserLogin)
	http.HandleFunc("/user/register",ctrl.UserRegister)
	http.HandleFunc("/contact/loadcommunity", ctrl.LoadCommunity)
	http.HandleFunc("/contact/loadfriend", ctrl.LoadFriend)
	http.HandleFunc("/contact/joincommunity", ctrl.JoinCommunity)
	http.HandleFunc("/contact/createcommunity", ctrl.CreateCommunity)
	http.HandleFunc("/contact/addfriend", ctrl.Addfriend)
	http.HandleFunc("/chat", ctrl.Chat)
	http.HandleFunc("/attach/upload", ctrl.Upload)
	//support the static resource
	http.Handle("/asset/", http.FileServer(http.Dir(".")))
	http.Handle("/mnt/", http.FileServer(http.Dir(".")))

	RegisterView()

	//http.HandleFunc("/user/login.shtml",
	//	func(writer http.ResponseWriter, request *http.Request) {
	//		tpl,err := template.ParseFiles("view/user/login.html")
	//		if err!=nil {
	//			//quit and print the err
	//			log.Fatal(err.Error())
	//		}
	//		tpl.ExecuteTemplate(writer,"/user/login.shtml",nil)
	//
	//})

	//start the web server
	http.ListenAndServe(":8080",nil)

}
