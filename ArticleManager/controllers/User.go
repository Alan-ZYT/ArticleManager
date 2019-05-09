package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ArticleManager/models"
	"encoding/base64"
)

type UserController struct {
	beego.Controller
}

//展示用户注册页面
func (this *UserController) ShowLogin() {
	//判定是否展现用户名
	userName := this.Ctx.GetCookie("UserName")
	uName, _ := base64.StdEncoding.DecodeString(userName)
	if userName != "" {
		this.Data["userName"] = string(uName)
		this.Data["checked"] = "checked"
	} else {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	}
	this.Layout = "userlayout.html"
	this.TplName = "login.html"
}

//处理用户登陆函数
func (this *UserController) HandleLogin() {
	userName := this.GetString("userName")
	passWord := this.GetString("password")
	if userName == "" || passWord == "" {
		beego.Info("用户名或密码有空值")
		return
	}
	//查询数据库用户是否存在
	o := orm.NewOrm()
	var user models.User
	user.UserName = userName
	o.Read(&user, "UserName")
	if user.PassWord != passWord {
		this.Redirect("/register", 302)
		return
	}
	//获取是否记住用户名
	remember := this.GetString("remember")
	//对汉字进行加密实现cookie存储
	uName := base64.StdEncoding.EncodeToString([]byte(userName))
	if remember == "on" {
		this.Ctx.SetCookie("UserName", uName, 6000)
	} else {
		this.Ctx.SetCookie("UserName", userName, -1)
	}
	this.SetSession("UserName", userName)

	this.Redirect("/article/index", 302)

}

//注册展示函数及处理函数
func (this *UserController) ShowRegister() {
	this.Layout = "userlayout.html"
	this.TplName = "register.html"
}

func (this *UserController) HandleRegister() {
	userName := this.GetString("userName")
	passWord := this.GetString("password")
	if userName == "" || passWord == "" {
		beego.Info("用户名或密码有空值")
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.UserName = userName
	user.PassWord = passWord
	_, err := o.Insert(&user)
	if err != nil {
		beego.Info("注册出错")
		return
	}
	this.Redirect("/login", 302)

}
