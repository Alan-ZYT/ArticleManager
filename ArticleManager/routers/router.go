package routers

import (
	"github.com/astaxie/beego"
	"ArticleManager/controllers"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/article/*", beego.BeforeExec, filterFunc)

	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandleRegister")
	//文章展示页面
	beego.Router("/article/index", &controllers.ArticleController{}, "get:ShowIndex")
	//文章类型显示页面
	beego.Router("/article/addtype", &controllers.ArticleController{}, "get:ShowAddType;post:HandleAddType")
	//删除添加的文章类型
	beego.Router("/article/deleteaddtype",&controllers.ArticleController{},"get:DeleteAddType")
	//退出登陆链接
	beego.Router("/article/logout", &controllers.ArticleController{}, "get:Logout")
	//添加文章页面
	beego.Router("/article/addarticle", &controllers.ArticleController{}, "get:ShowAddArticle;post:AddArticle")
	//文章详情查看页面
	beego.Router("/article/content", &controllers.ArticleController{}, "get:ShowContent")
	//文章删除
	beego.Router("/article/delete", &controllers.ArticleController{}, "get:DeleteArticle")
	beego.Router("/article/update", &controllers.ArticleController{}, "get:ShowUpdateArticle;post:UpdateArticle")
}

func filterFunc(ctx *context.Context) {
	userName := ctx.Input.Session("UserName")
	if userName == nil {
		ctx.Redirect(302, "/login")
	}
}
