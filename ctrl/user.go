package ctrl

import (
	"fmt"
	"im/model"
	"im/service"
	"im/util"
	"log"
	"math/rand"
	"net/http"
)

var userService service.UserService

func UserRegister(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()

	mobile := request.PostForm.Get("mobile")
	plainpwd := request.PostForm.Get("passwd")

	nickname := fmt.Sprintf("user%06d", rand.Int31())
	avatar := ""
	sex := model.SEX_UNKNOW

	user, err := userService.Register(mobile, plainpwd, nickname, avatar, sex)
	if err!=nil {
		util.RespFail(writer,err.Error())
	}else {
		util.RespOk(writer,user,"")
	}
}

func UserLogin(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()
	mobile := request.PostForm.Get("mobile")
	passwd := request.PostForm.Get("passwd")

	//check login
	user, err := userService.Login(mobile,passwd)

	if err!=nil {
		util.RespFail(writer, err.Error())
	}else {
		util.RespOk(writer, user, "")
		log.Println("---------login----------")
		log.Println(user.Id)
		log.Println(user.Token)
	}
}
