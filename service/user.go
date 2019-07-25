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

		tmp := model.User{}
		//find the user by mobile
		DbEngin.Where("mobile = ?", mobile).Get(&tmp)
		//can not find the user
		if tmp.Id == 0 {
			return tmp,errors.New("no user")
		}

		//check the pwd
		if !util.ValidatePasswd(plainpwd,tmp.Salt,tmp.Passwd){
			return tmp,errors.New("password wrong")
		}

		//flush token
		str := fmt.Sprintf("%d", time.Now().Unix())
		token := util.MD5Encode(str)
		tmp.Token = token
		DbEngin.Id(tmp.Id).Cols("token").Update(&tmp)


		return user,nil
}
func (s *UserService)Find(
	userId int64 )(user model.User) {

	//首先通过手机号查询用户
	tmp :=model.User{

	}
	DbEngin.ID(userId).Get(&tmp)
	return tmp
}