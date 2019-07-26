# IM系统架构

毕设项目接触了IM，当时是用Java实现，现在来用golang复现一下，有时间再复盘整理一下毕设项目

源代码：<https://github.com/rexllz/im>

![i1](<https://raw.githubusercontent.com/rexllz/im/master/img/i1.jpg>)

<!--more-->

# 单机性能瓶颈

- Map

Map不能太大

读写锁（读次数远大于写次数）

- 系统

Linux的最大文件数影响

- CPU

JSON编码次数影响最大

IO资源的使用（合并写操作）

多使用缓存

- 应用/资源服务相分离

文件服务迁移到oss

# 搭建框架

## 前端获取数据DEMO

```go
package main

import (
	"net/http"
)

func main() {

	//bind the func and request
	http.HandleFunc("/user/login",
		func(writer http.ResponseWriter, request *http.Request) {

			request.ParseForm()
			mobile := request.PostForm.Get("mobile")
			passwd := request.PostForm.Get("passwd")
			loginok := false

			if (mobile == "186" && passwd == "123") {

				loginok = true
			}

			str := `{"code":0,"data":{"id":1,"token":"test"}}`

			if !loginok {

				str = `{"code":-1,"msg":"password wrong"}`
			}
			//set header
			writer.Header().Set("Content-Type","application/json")
			writer.WriteHeader(http.StatusOK)

			writer.Write([]byte(str))

			//io.WriteString(writer,"hello world")
	})

	//start the web server
	http.ListenAndServe(":8080",nil)
}
```

![i0](https://raw.githubusercontent.com/rexllz/im/master/img/i0.jpg)

## View的使用和支持

```go
//support the static resource
	http.Handle("/", http.FileServer(http.Dir(".")))
```

create the shtml

```html
{{define "/user/login.shtml"}}
<!DOCTYPE html>
<html>
</html>
<script>
</script>
{{end}}
```

handle

```go
http.HandleFunc("/user/login.shtml",
		func(writer http.ResponseWriter, request *http.Request) {
			tpl,err := template.ParseFiles("view/user/login.html")
			if err!=nil {
				//quit and print the err
				log.Fatal(err.Error())
			}
			tpl.ExecuteTemplate(writer,"/user/login.shtml",nil)

	})
```

![i3](https://raw.githubusercontent.com/rexllz/im/master/img/i3.jpg)

## Xorm操作数据库

```
go get github.com/go-xorm/xorm
go get github.com/go-sql-driver/mysql
```

定义init函数（自动运行）

```go
var DbEngin *xorm.Engine
func init(){
	drivename := "mysql"
	DsName := "root:root@(127.0.0.1:3306)/imchat?charset=utf8"
	DbEngin, err := xorm.NewEngine(drivename,DsName)
	if err!=nil {
		log.Fatal(err.Error())
	}
	//show the sql
	DbEngin.ShowSQL(true)
	//set the max connect num
	DbEngin.SetMaxOpenConns(2)
	//auto create tables
	//DbEngin.Sync2(new(User))
	fmt.Println("init DB connect")
}
```

## 创建web项目结构

- 建立目录
- init.go（DB初始化）
- 服务函数

实现注册功能

service

```go
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
```

model

```go
package model

import "time"

const (
	SEX_WOMEN="W"
	SEX_MEN="M"
	//
	SEX_UNKNOW="U"
)
type User struct {
	//用户ID
	Id         int64     `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`
	Mobile   string 		`xorm:"varchar(20)" form:"mobile" json:"mobile"`
	Passwd       string	`xorm:"varchar(40)" form:"passwd" json:"-"`   // 什么角色
	Avatar	   string 		`xorm:"varchar(150)" form:"avatar" json:"avatar"`
	Sex        string	`xorm:"varchar(2)" form:"sex" json:"sex"`   // 什么角色
	Nickname    string	`xorm:"varchar(20)" form:"nickname" json:"nickname"`   // 什么角色
	//加盐随机字符串6
	Salt       string	`xorm:"varchar(10)" form:"salt" json:"-"`   // 什么角色
	Online     int	`xorm:"int(10)" form:"online" json:"online"`   //是否在线
	//前端鉴权因子,
	Token      string	`xorm:"varchar(40)" form:"token" json:"token"`   // 什么角色
	Memo      string	`xorm:"varchar(140)" form:"memo" json:"memo"`   // 什么角色
	Createat   time.Time	`xorm:"datetime" form:"createat" json:"createat"`   // 什么角色
}
```

ctrl

```go
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

```

util

```go
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

```

![i6](https://raw.githubusercontent.com/rexllz/im/master/img/i6.jpg)

![i5](https://raw.githubusercontent.com/rexllz/im/master/img/i5.jpg)

# IM业务

websocket.connect为并行数据，需要转为串行

```go
node := &Node{
		Conn:conn,

		//transfer Parallel to serial
		DataQueue:make(chan []byte,50),
		GroupSets:set.New(set.ThreadSafe),
	}

```

存储node，user的map需要上读写锁，以防止出错

读写锁实际是一种特殊的[自旋锁](https://baike.baidu.com/item/自旋锁)，它把对共享资源的访问者划分成读者和写者，读者只对共享资源进行读访问，写者则需要对共享资源进行写操作。这种锁相对于自旋锁而言，能提高并发性，因为在[多处理器系统](https://baike.baidu.com/item/多处理器系统)中，它允许同时有多个读者来访问共享资源，最大可能的读者数为实际的逻辑CPU数。写者是排他性的，一个读写锁同时只能有一个写者或多个读者（与CPU数相关），但不能同时既有读者又有写者。

```go
rwlocker.Lock()
	clientMap[userId]=node
	rwlocker.Unlock()

```



```go
//todo 完成发送逻辑,con
	go sendproc(node)
	//todo 完成接收逻辑
	go recvproc(node)

```



设计可以无限扩张业务场景的消息通讯结构

```cgo
func recvproc(node *Node) {
	for{
		_,data,err := node.Conn.ReadMessage()
		if err!=nil{
			log.Println(err.Error())
			return
		}
		//todo 对data进一步处理
		//dispatch(data)
		fmt.Printf("recv<=%s",data)
	}
}

```

## 原理

前端通过websocket发送`json格式的字符串`
用户2向用户3发送文字消息hello

```json5
{id:1,userid:2,dstid:3,cmd:10,media:1,content:"hello"}

```

里面携带
谁发的-userid
要发给谁-dstid
这个消息有什么用-cmd
消息怎么展示-media
消息内容是什么-(url,amout,pic,content等)

## 核心数据结构

```go
type Message struct {
	Id      int64  `json:"id,omitempty" form:"id"` //消息ID
	//谁发的
	Userid  int64  `json:"userid,omitempty" form:"userid"` //谁发的
	//什么业务
	Cmd     int    `json:"cmd,omitempty" form:"cmd"` //群聊还是私聊
	//发给谁
	Dstid   int64  `json:"dstid,omitempty" form:"dstid"`//对端用户ID/群ID
	//怎么展示
	Media   int    `json:"media,omitempty" form:"media"` //消息按照什么样式展示
	//内容是什么
	Content string `json:"content,omitempty" form:"content"` //消息的内容
	//图片是什么
	Pic     string `json:"pic,omitempty" form:"pic"` //预览图片
	//连接是什么
	Url     string `json:"url,omitempty" form:"url"` //服务的URL
	//简单描述
	Memo    string `json:"memo,omitempty" form:"memo"` //简单描述
	//其他的附加数据，语音长度/红包金额
	Amount  int    `json:"amount,omitempty" form:"amount"` //其他和数字相关的
}
const (
    //点对点单聊,dstid是用户ID
	CMD_SINGLE_MSG = 10
	//群聊消息,dstid是群id
	CMD_ROOM_MSG   = 11
	//心跳消息,不处理
	CMD_HEART      = 0
	
)
const (
    //文本样式
	MEDIA_TYPE_TEXT=1
	//新闻样式,类比图文消息
	MEDIA_TYPE_News=2
	//语音样式
	MEDIA_TYPE_VOICE=3
	//图片样式
	MEDIA_TYPE_IMG=4
	
	//红包样式
	MEDIA_TYPE_REDPACKAGR=5
	//emoj表情样式
	MEDIA_TYPE_EMOJ=6
	//超链接样式
	MEDIA_TYPE_LINK=7
	//视频样式
	MEDIA_TYPE_VIDEO=8
	//名片样式
	MEDIA_TYPE_CONCAT=9
	//其他自己定义,前端做相应解析即可
	MEDIA_TYPE_UDEF=100
)
/**
消息发送结构体,点对点单聊为例
1、MEDIA_TYPE_TEXT
{id:1,userid:2,dstid:3,cmd:10,media:1,
content:"hello"}

3、MEDIA_TYPE_VOICE,amount单位秒
{id:1,userid:2,dstid:3,cmd:10,media:3,
url:"http://www.a,com/dsturl.mp3",
amount:40}

4、MEDIA_TYPE_IMG
{id:1,userid:2,dstid:3,cmd:10,media:4,
url:"http://www.baidu.com/a/log.jpg"}


2、MEDIA_TYPE_News
{id:1,userid:2,dstid:3,cmd:10,media:2,
content:"标题",
pic:"http://www.baidu.com/a/log,jpg",
url:"http://www.a,com/dsturl",
"memo":"这是描述"}


5、MEDIA_TYPE_REDPACKAGR //红包amount 单位分
{id:1,userid:2,dstid:3,cmd:10,media:5,url:"http://www.baidu.com/a/b/c/redpackageaddress?id=100000","amount":300,"memo":"恭喜发财"}
6、MEDIA_TYPE_EMOJ 6
{id:1,userid:2,dstid:3,cmd:10,media:6,"content":"cry"}

7、MEDIA_TYPE_Link 7
{id:1,userid:2,dstid:3,cmd:10,media:7,
"url":"http://www.a.com/dsturl.html"
}

8、MEDIA_TYPE_VIDEO 8
{id:1,userid:2,dstid:3,cmd:10,media:8,
pic:"http://www.baidu.com/a/log,jpg",
url:"http://www.a,com/a.mp4"
}

9、MEDIA_TYPE_CONTACT 9
{id:1,userid:2,dstid:3,cmd:10,media:9,
"content":"10086",
"pic":"http://www.baidu.com/a/avatar,jpg",
"memo":"胡大力"}

*/

```

从哪里接收数据?怎么处理这些数据呢?

```go
func recvproc(node *Node) {
	for{
		_,data,err := node.Conn.ReadMessage()
		if err!=nil{
			log.Println(err.Error())
			return
		}
		//todo 对data进一步处理
		fmt.Printf("recv<=%s",data)
		dispatch(data)
	}
}
func dispatch(data []byte){
    //todo 转成message对象
    
    //todo 根据cmd参数处理逻辑
    
    
    
    
    
    msg :=Message{}
    err := json.UnMarshal(data,&msg)
    if err!=nil{
        log.Printf(err.Error())
        return ;
    }
    switch msg.Cmd {
    	case CMD_SINGLE_MSG: //如果是单对单消息,直接将消息转发出去
    		//向某个用户发回去
    		fmt.Printf("c2cmsg %d=>%d\n%s\n",msg.Userid,msg.Dstid,string(tmp))
    		SendMsgToUser(msg.Userid, msg.Dstid, tmp)
    		//fmt.Println(msg)
    	case CMD_ROOM_MSG: //群聊消息,需要知道
    		fmt.Printf("c2gmsg %d=>%d\n%s\n",msg.Userid,msg.Dstid,string(tmp))
    		SendMsgToRoom(msg.Userid, msg.Dstid, tmp)
    	case CMD_HEART:
    	default:
    	    //啥也别做
    	    
    	}
    		
}

```

## 发送文字、表情包

前端user1拼接好数据对象Message
msg={id:1,userid:2,dstid:3,cmd:10,media:1,content:txt}
转化成json字符串jsonstr
jsonstr = JSON.stringify(msg)
通过websocket.send(jsonstr)发送
后端S在recvproc中接收收数据data
并做相应的逻辑处理dispatch(data)-转发给user2
user2通过websocket.onmessage收到消息后做解析并显示

前端所有的操作都在拼接数据
如何拼接?

```javascript
sendtxtmsg:function(txt){
//{id:1,userid:2,dstid:3,cmd:10,media:1,content:txt}
var msg =this.createmsgcontext();
//msg={"dstid":dstid,"cmd":cmd,"userid":userId()}
//选择某个好友/群的时候对dstid,cmd进行赋值
//userId()返回用户自己的id ,
// 从/chat/index.shtml?id=xx&token=yy中获得
//1文本类型
msg.media=1;msg.content=txt;
this.showmsg(userInfo(),msg);//显示自己发的文字
this.webSocket.send(JSON.stringify(msg))//发送
}

sendpicmsg:function(picurl){
    //{id:1,userid:2,dstid:3,cmd:10,media:4,
    // url:"http://www.baidu.com/a/log,jpg"}
    var msg =this.createmsgcontext();
    msg.media=4;
    msg.url=picurl;
    this.showmsg(userInfo(),msg)
    this.webSocket.send(JSON.stringify(msg))
}
sendaudiomsg:function(url,num){
    //{id:1,userid:2,dstid:3,cmd:10,media:3,url:"http://www.a,com/dsturl.mp3",anount:40}
    var msg =this.createmsgcontext();
    msg.media=3;
    msg.url=url;
    msg.amount = num;
    this.showmsg(userInfo(),msg)
    console.log("sendaudiomsg",this.msglist);
    this.webSocket.send(JSON.stringify(msg))
}

```

### 后端逻辑处理函数 

func dispatch(data[]byte)

```cgo
func dispatch(data[]byte){
    //todo 解析data为message
    
    //todo根据message的cmd属性做相应的处理
    
}
func recvproc(node *Node) {
	for{
		_,data,err := node.Conn.ReadMessage()
		if err!=nil{
			log.Println(err.Error())
			return
		}
		//todo 对data进一步处理
		dispatch(data)
		fmt.Printf("recv<=%s",data)
	}
}

```

### 对端接收到消息后处理函数

```js
//初始化websocket的时候进行回调配置
this.webSocket.onmessage = function(evt){
     //{"data":"}",...}
     if(evt.data.indexOf("}")>-1){
         this.onmessage(JSON.parse(evt.data));
     }else{
         console.log("recv<=="+evt.data)
     }
 }.bind(this)
onmessage:function(data){
     this.loaduserinfo(data.userid,function(user){
         this.showmsg(user,data)
     }.bind(this))
 }

 //消息显示函数
showmsg:function(user,msg){
    var data={}
    data.ismine = userId()==msg.userid;
    //console.log(data.ismine,userId(),msg.userid)
    data.user = user;
    data.msg = msg;
    //vue 只需要修改数据结构即可完成页面渲染
    this.msglist = this.msglist.concat(data)
    //面板重置
    this.reset();
    var that =this;
    //滚动到新消息处
    that.timer = setTimeout(function(){
        window.scrollTo(0, document.getElementById("convo").offsetHeight);
        clearTimeout(that.timer)
    },100)
 }

```

### 表情包简单逻辑

弹出一个窗口,
选择图片获得一个连接地址
调用sendpicmsg方法开始发送流程

## 图片等upload

##5.5 发送图片/拍照
弹出一个窗口,
选择图片,上传到服务器
获得一个链接地址
调用sendpicmsg方法开始发送流程

```html
<input 
accept="image/gif,image/jpeg,,image/png" 
type="file" 
onchange="upload(this)" 
class='upload'/>

```

sendpicmsg方法开始发送流程

### upload前端实现

```javascript
function upload(dom){
        uploadfile("attach/upload",dom,function(res){
            if(res.code==0){//成功以后调用sendpicmsg
                vm.sendpicmsg(res.data)
            }
        })
    }
    
function uploadfile(uri,dom,callback){
    //H5新特性
    var formdata = new FormData();
    //获得一个文件dom.files[0]
    formdata.append("file",dom.files[0])
    //formdata.append("filetype",".png")//.mp3指定后缀
    
    var xhr = new XMLHttpRequest();//ajax初始化
    var url = "http://"+location.host+"/"+uri;
    //"http://127.0.0.1/attach/upload"
    xhr.open("POST",url, true);
    //成功时候回调
    xhr.onreadystatechange = function() {
        if (xhr.readyState == 4 && 
        xhr.status == 200) {
            //fn.call(this, JSON.parse(xhr.responseText));
            callback(JSON.parse(xhr.responseText))
        }
    };
    xhr.send(formdata);
}    

```

### upload后端实现

存储到本地

```
func UploadLocal(writer http.ResponseWriter,
	request * http.Request){
	}

```

存储到alioss

```
func UploadLocal(writer http.ResponseWriter,
	 request * http.Request){
}
如何安装 golang.org/x/time/rate
>cd $GOPATH/src/golang.org/x/
>git clone https://github.com/golang/time.git time


```

采集语音

```javascript
navigator.mediaDevices.getUserMedia(
    {audio: true, video: true}
    ).then(successfunc).catch(errfunc);


navigator.mediaDevices.getUserMedia(
    {audio: true, video: false}
    ).then(function(stream)  {
              //请求成功
              this.recorder = new MediaRecorder(stream);
              this.recorder.start();
              this.recorder.ondataavailable = (event) => {
                  uploadblob("attach/upload",event.data,".mp3",res=>{
                      var duration = Math.ceil((new Date().getTime()-this.duration)/1000);
                      this.sendaudiomsg(res.data,duration);
                  })

                  stream.getTracks().forEach(function (track) {
                      track.stop();
                  });
                  this.showprocess = false
              }
              
          }.bind(this)).catch(function(err){
                mui.toast(err.msg)
                this.showprocess = false
            }.bind(this));

```

上传语音

```javascript
function uploadblob(uri,blob,filetype,fn){
       var xhr = new XMLHttpRequest();
       xhr.open("POST","//"+location.host+"/"+uri, true);
       // 添加http头，发送信息至服务器时内容编码类型
       xhr.onreadystatechange = function() {
           if (xhr.readyState == 4 && (xhr.status == 200 || xhr.status == 304)) {
               fn.call(this, JSON.parse(xhr.responseText));
           }
       };
       var _data=[];
       var formdata = new FormData();
       formdata.append("filetype",filetype);
       formdata.append("file",blob)
       xhr.send(formdata);
   }

```

## 群聊

分析群id,找到加了这个群的用户,把消息发送过去

### 方案一

map<userid><qunid1,qunid2,qunid3>
优势是锁的频次低
劣势是要轮训全部map

```go
type Node struct {
	Conn *websocket.Conn
	//并行转串行,
	DataQueue chan []byte
	GroupSets set.Interface
}
//映射关系表
var clientMap map[int64]*Node = make(map[int64]*Node,0)

```

### 方案二

map<群id><userid1,userid2,userid3>
优势是找用户ID非常快
劣势是发送信息时需要根据userid获取node,锁的频次太高

```go
type Node struct {
	Conn *websocket.Conn
	//并行转串行,
	DataQueue chan []byte
}
//映射关系表
var clientMap map[int64]*Node = make(map[int64]*Node,0)
var comMap map[int64]set.Interface= make(map[int64]set.Interface,0)


```

需要处理的问题

```javascript
1、当用户接入的时候初始化groupset
2、当用户加入群的时候刷新groupset
3、完成信息分发

```

# 优化

## 静态资源分离（Aliyun OSS）

使用阿里云OSS

> go get github.com/aliyun/aliyun-oss-go-sdk/oss

```go
//权限设置为公共读状态
//需要安装
func UploadOss(writer http.ResponseWriter,
	request * http.Request){
	//todo 获得上传的文件
	srcfile,head,err:=request.FormFile("file")
	if err!=nil{
		util.RespFail(writer,err.Error())
		return
	}


	//todo 获得文件后缀.png/.mp3

	suffix := ".png"
	//如果前端文件名称包含后缀 xx.xx.png
	ofilename := head.Filename
	tmp := strings.Split(ofilename,".")
	if len(tmp)>1{
		suffix = "."+tmp[len(tmp)-1]
	}
	//如果前端指定filetype
	//formdata.append("filetype",".png")
	filetype := request.FormValue("filetype")
	if len(filetype)>0{
		suffix = filetype
	}

	//todo 初始化ossclient
	client,err:=oss.New(EndPoint,AccessKeyId,AccessKeySecret)
	if err!=nil{
		util.RespFail(writer,err.Error())
		return
	}
	//todo 获得bucket
	bucket,err := client.Bucket(Bucket)
	if err!=nil{
		util.RespFail(writer,err.Error())
		return
	}
	//todo 设置文件名称
	//time.Now().Unix()
	filename := fmt.Sprintf("mnt/%d%04d%s",
		time.Now().Unix(), rand.Int31(),
		suffix)
	//todo 通过bucket上传
	err=bucket.PutObject(filename,srcfile)
	if err!=nil{
		util.RespFail(writer,err.Error())
		return
	}
	//todo 获得url地址
	url := "http://"+Bucket+"."+EndPoint+"/"+filename

	//todo 响应到前端
	util.RespOk(writer,url,"")
}

```

# 分布式部署

## nginx反向代理

![i7](https://raw.githubusercontent.com/rexllz/im/master/img/i7.jpg)

普通方案无法满足connect之间无法通讯的问题

需要建立总线维持信息

## 消息总线

![i8](https://raw.githubusercontent.com/rexllz/im/master/img/i8.jpg)

## 局域网通讯协议

![i9](https://raw.githubusercontent.com/rexllz/im/master/img/i9.jpg)

## 实现调度应用

![i10](https://raw.githubusercontent.com/rexllz/im/master/img/i10.jpg)

## 实现（UDP方案）

回顾单体应用
开启ws接收协程recvproc/ws发送协程sendproc
websocket收到消息->dispatch发送给dstid

基于UDP的分布式应用
开启ws接收协程recvproc/ws发送协程sendproc
开启udp接收协程udprecvproc/udp发送协程udpsendproc

websocket收到消息->broadMsg广播到局域网
udp接收到收到消息->dispatch发送给dstid
自己是局域网一份子,所以也能接收到消息

```go
var  udpsendchan chan []byte=make(chan []byte,1024)
//todo 将消息广播到局域网
func broadMsg(data []byte){
	udpsendchan<-data
}

//todo 完成udp数据的发送协程
func udpsendproc(){
	log.Println("start udpsendproc")
	//todo 使用udp协议拨号
	con,err:=net.DialUDP("udp",nil,
		&net.UDPAddr{
			IP:net.IPv4(192,168,0,255),
			Port:3000,
		})
	defer con.Close()
	if err!=nil{
		log.Println(err.Error())
		return
	}
	//todo 通过的到的con发送消息
	//con.Write()
	for{
		select {
		case data := <- udpsendchan:
			_,err=con.Write(data)
			if err!=nil{
				log.Println(err.Error())
				return
			}
		}
	}
}
//todo 完成upd接收并处理功能
func udprecvproc(){
	log.Println("start udprecvproc")
	//todo 监听udp广播端口
	con,err:=net.ListenUDP("udp",&net.UDPAddr{
		IP:net.IPv4zero,
		Port:3000,
	})
	defer con.Close()
	if err!=nil{log.Println(err.Error())}
	//TODO 处理端口发过来的数据
	for{
		var buf [512]byte
		n,err:=con.Read(buf[0:])
		if err!=nil{
			log.Println(err.Error())
			return
		}
		//直接数据处理
		dispatch(buf[0:n])
	}
	log.Println("stop updrecvproc")
}

```



### nginx反向代理

```
	upstream wsbackend {
			server 192.168.0.102:8080;
			server 192.168.0.100:8080;
			hash $request_uri;
	}
	map $http_upgrade $connection_upgrade {
    default upgrade;
    ''      close;
	}
    server {
	  listen  80;
	  server_name localhost;
	  location / {
	   proxy_pass http://wsbackend;
	  }
	  location ^~ /chat {
	   proxy_pass http://wsbackend;
	   proxy_connect_timeout 500s;
       proxy_read_timeout 500s;
	   proxy_send_timeout 500s;
	   proxy_set_header Upgrade $http_upgrade;
       proxy_set_header Connection "Upgrade";
	  }
	 }

}

```

### 打包发布

- windows平台

```bash
::remove dir
rd /s/q release
::make dir 
md release
::go build -ldflags "-H windowsgui" -o chat.exe
go build -o chat.exe
::
COPY chat.exe release\
COPY favicon.ico release\favicon.ico
::
XCOPY asset\*.* release\asset\  /s /e
XCOPY view\*.* release\view\  /s /e 

```

- linux平台

```bash
#!/bin/sh
rm -rf ./release
mkdir  release
go build -o chat
chmod +x ./chat
cp chat ./release/
cp favicon.ico ./release/
cp -arf ./asset ./release/
cp -arf ./view ./release/

```

- 运行注意事项
  linux 下

```bash
nohup ./chat >>./log.log 2>&1 &

```

# Tips

## idea .ignore插件

![i2](https://raw.githubusercontent.com/rexllz/im/master/img/i2.jpg)

## JSON改为小写

nil值不发送

```go
type H struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

```

## 自动渲染和接入全部View

写入统一函数

```go
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

```

## 自动创建表结构

```go
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


```

![i4](https://raw.githubusercontent.com/rexllz/im/master/img/i4.jpg)

