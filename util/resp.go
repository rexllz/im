package util

import (
	"encoding/json"
	"log"
	"net/http"
)


//type H struct {
//	Code int `json:"code"`
//	Msg string `json:"msg"`
//	Data interface{} `json:"data,omitempty"`
//}
type H struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
	Data interface{} `json:"data,omitempty"`
	Rows interface{} `json:"rows,omitempty"`
	Total interface{} `json:"total,omitempty"`
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
func RespList(w http.ResponseWriter,code int,data interface{},total interface{})  {

	w.Header().Set("Content-Type","application/json")
	//设置200状态
	w.WriteHeader(http.StatusOK)
	//输出
	//定义一个结构体
	//满足某一条件的全部记录数目
	//测试 100
	//20
	h := H{
		Code:code,
		Rows:data,
		Total:total,
	}
	//将结构体转化成JSOn字符串
	ret,err := json.Marshal(h)
	if err!=nil{
		log.Println(err.Error())
	}
	//输出
	w.Write(ret)
}
func RespOkList(w http.ResponseWriter,lists interface{},total interface{}){
	//分页数目,
	RespList(w,0,lists,total)
}