package util

import (
	"encoding/json"
	"log"
	"net/http"
)

type H struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func RespFail(w http.ResponseWriter,msg string){
	Resp(w,-1,nil,msg)
}

func RespOk(w http.ResponseWriter,data interface{},msg string){
	Resp(w,0,data,msg)
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