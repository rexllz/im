package service

import "../model"

type UserService struct {

}

func (s* UserService)Register (
	mobile,
	plainpwd,
	nickname,
	avatar,
	sex string)(user model.User, err error){
	return user,nil
}

func (s* UserService)Login (
	mobile,
	plainpwd string)(user model.User, err error){
	return user,nil
}