package service

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"im/model"
	"log"
)

var DbEngin *xorm.Engine
func init(){
	drivename := "mysql"
	DsName := "root:root@(127.0.0.1:3306)/imchat?charset=utf8"
	err := errors.New("")
	DbEngin, err = xorm.NewEngine(drivename,DsName)
	if err!=nil && ""!=err.Error(){
		log.Fatal(err.Error())
	}
	//show the sql
	DbEngin.ShowSQL(true)
	//set the max connect num
	DbEngin.SetMaxOpenConns(2)
	//auto create tables
	DbEngin.Sync2(new(model.User))
	fmt.Println("init DB connect")
}
