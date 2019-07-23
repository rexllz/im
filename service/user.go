package service

import (
	"errors"
	"fmt"
	"im/model"
	"im/util"
	"math/rand"
	"time"
)

type UserService struct {

}

func (s* UserService)Register (
	mobile,
	plainpwd,
	nickname,
	avatar,
	sex string)(user model.User, err error){

		//check if the mobile exist
		tmp := model.User{}
	    _, err = DbEngin.Where("mobile=?", mobile).Get(&tmp)
		if err!=nil {
			return tmp,err
		}
		if tmp.Id>0 {
			return tmp,errors.New("this mobile have account")
		}

		tmp.Mobile = mobile
		tmp.Nickname = nickname
		tmp.Avatar = avatar
		tmp.Sex = sex
		tmp.Salt = fmt.Sprintf("%06d",rand.Int31n(10000))
		tmp.Passwd = util.MakePasswd(plainpwd,tmp.Salt)
		tmp.Createat = time.Now()
		tmp.Token = fmt.Sprintf("%08d",rand.Int31())

		//insert one data,return number and error
		_, err = DbEngin.InsertOne(&tmp)


		return tmp, err
}

func (s* UserService)Login (
	mobile,
	plainpwd string)(user model.User, err error){
	return user,nil
}