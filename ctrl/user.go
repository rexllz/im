package ctrl

import (
	"fmt"
	"im/model"
	"im/service"
	"im/util"
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
	loginok := false

	if (mobile == "186" && passwd == "123") {

		loginok = true
	}

	if loginok {
		data := make(map[string]interface{})
		data["id"] = 1
		data["token"] = "test"
		util.RespOk(writer,data,"")
	}else {
		util.RespFail(writer,"password wrong")
	}

}
