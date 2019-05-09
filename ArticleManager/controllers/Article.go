package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ArticleManager/models"
	"path"
	"time"
	"math"
	"strconv"
	"github.com/gomodule/redigo/redis"
	"bytes"
	"encoding/gob"
)

type ArticleController struct {
	beego.Controller
}

/*************************文章列表页面******************************/
//文章列表展示页面
func (this *ArticleController) ShowIndex() {
	Typename := this.GetString("select")
	this.Data["TypeName"] = Typename
	//下拉框展示
	o := orm.NewOrm()
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	defer conn.Close()
	if err != nil {
		beego.Info("conn err:", err)
		return
	}
	reply, err := conn.Do("get", "articletype")
	b, _ := redis.Bytes(reply, err)
	//beego.Info(b)
	var ArticleType []models.ArticleType
	if len(b) == 0 {
		//beego.Info(111111)
		o.QueryTable("ArticleType").All(&ArticleType)
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		enc.Encode(ArticleType)
		conn.Do("set", "articletype", buffer.Bytes())
	} else {
		dec:=gob.NewDecoder(bytes.NewBuffer(b))
		dec.Decode(&ArticleType)
		beego.Info(ArticleType)

	}
	this.Data["itemList"] = ArticleType
	//设置文章记录数目d
	var RecordCount int64
	RecordCount, _ = o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__TypeName", Typename).Count()

	this.Data["RecordCount"] = RecordCount
	//设置文章分页
	EveryPageCount := 1
	PageNum := math.Ceil(float64(RecordCount) / float64(EveryPageCount))

	this.Data["PageNum"] = PageNum //设置总页数
	PageIndex, _ := this.GetInt("PageIndex")
	this.Data["Index"] = PageIndex
	//加载文章列表,并根据下拉框的选中类型展示列表
	var ArticleList []models.Article
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__TypeName", Typename).Limit(EveryPageCount, EveryPageCount*(PageIndex-1)).All(&ArticleList)
	this.Data["ArticleList"] = ArticleList

	//页面加载数据
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Scripts"] = "scripts.html"
	this.Layout = "articlelayout.html"
	this.TplName = "index.html"
	userName := this.GetSession("UserName")
	this.Data["userName"] = userName
}

//退出登陆设置
func (this *ArticleController) Logout() {
	this.DelSession("UserName")
	this.Redirect("/login", 302)
}

func (this *ArticleController) ShowContent() {
	userName := this.GetSession("UserName")
	this.Data["userName"] = userName
	ArtitleTitle := this.GetString("title")
	o := orm.NewOrm()
	var Article models.Article
	Article.ArticleType = new(models.ArticleType)
	Article.Title = ArtitleTitle
	o.Read(&Article, "Title")
	var Article1 models.Article
	o.QueryTable(Article).RelatedSel("ArticleType").Filter("ArticleType__Id", Article.Id).One(&Article1)
	this.Data["Article"] = Article1
	Article.ReadCount += 1
	o.Update(&Article)

	//展示用户浏览记录
	var Users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id", Article.Id).Distinct().All(&Users)
	this.Data["Users"] = Users

	//添加用户浏览记录
	m2m := o.QueryM2M(&Article, "Users")
	User := models.User{UserName: userName.(string)}
	o.Read(&User, "UserName")

	m2m.Add(User)
	this.Layout = "articlelayout.html"
	this.TplName = "content.html"
}
func (this *ArticleController) DeleteArticle() {
	id, _ := this.GetInt("id")
	var Article models.Article
	Article.Id = id
	o := orm.NewOrm()
	o.Read(&Article)
	o.Delete(&Article)
	this.Redirect("/article/index", 302)
}

func (this *ArticleController) ShowUpdateArticle() {
	userName := this.GetSession("UserName")
	this.Data["userName"] = userName
	Id, _ := this.GetInt("id")
	beego.Info(Id)
	o := orm.NewOrm()
	var Article models.Article
	Article.Id = Id
	o.Read(&Article)
	this.Data["Article"] = Article

	this.Layout = "articlelayout.html"
	this.TplName = "update.html"
}
func (this *ArticleController) UpdateArticle() {
	Titile := this.GetString("articleName")
	Content := this.GetString("content")
	id, _ := this.GetInt("id")
	if Titile == "" {
		this.Redirect("/article/update?id="+strconv.Itoa(id), 302)
		return
	}
	if Content == "" {
		this.Redirect("/article/update?id="+strconv.Itoa(id), 302)
		return
	}
	file, head, err := this.GetFile("uploadname")
	if err != nil {
		this.Redirect("/article/update?id="+strconv.Itoa(id), 302)
	}
	defer file.Close()
	if head.Size > 500000 {
		beego.Info("上传图片尺寸过大")
		this.Redirect("/article/update?id="+strconv.Itoa(id), 302)
		return
	}
	if path.Ext(head.Filename) != ".jpg" && path.Ext(head.Filename) != ".jpeg" && path.Ext(head.Filename) != ".png" {
		beego.Info("上传文件不是图片")
		this.Redirect("/article/update?id="+strconv.Itoa(id), 302)
		return
	}
	//为上传文件更名
	filename := time.Now().Format("20060102150405") + path.Ext(head.Filename)
	this.SaveToFile("uploadname", "./static/img/"+filename)
	o := orm.NewOrm()
	var Article models.Article
	Article.Title = Titile
	Article.Content = Content
	Article.Img = "/static/img/" + filename
	o.Insert(&Article)
	this.Redirect("/article/index", 302)
}

/**************************添加文章页面********************************/
//添加文章页面
func (this *ArticleController) ShowAddArticle() {
	//从数据库中获取文章类型列表
	o := orm.NewOrm()
	var TypeList []models.ArticleType
	o.QueryTable("ArticleType").All(&TypeList)
	this.Data["itemlist"] = TypeList
	this.Layout = "articlelayout.html"
	this.TplName = "add.html"
	userName := this.GetSession("UserName")
	this.Data["userName"] = userName

}

func (this *ArticleController) AddArticle() {
	Titile := this.GetString("articleName")
	TypeName := this.GetString("select")
	Content := this.GetString("content")
	if Titile == "" {
		this.Redirect("/article/addarticle", 302)
		return
	}
	if Content == "" {
		this.Redirect("/article/addarticle", 302)
		return
	}
	file, head, err := this.GetFile("uploadname")
	if err != nil {
		this.Redirect("/article/addarticle", 302)
	}
	defer file.Close()
	if head.Size > 500000 {
		beego.Info("上传图片尺寸过大")
		this.Redirect("/article/addarticle", 302)
		return
	}
	if path.Ext(head.Filename) != ".jpg" && path.Ext(head.Filename) != ".jpeg" && path.Ext(head.Filename) != ".png" {
		beego.Info("上传文件不是图片")
		this.Redirect("/article/addarticle", 302)
		return
	}
	//为上传文件更名
	filename := time.Now().Format("20060102150405") + path.Ext(head.Filename)
	this.SaveToFile("uploadname", "./static/img/"+filename)

	//将数据插入数据库
	o := orm.NewOrm()
	var Article models.Article
	var Articletype = new(models.ArticleType)
	Articletype.TypeName = TypeName
	o.Read(Articletype, "TypeName")
	Article.Title = Titile
	Article.ArticleType = Articletype
	Article.Content = Content
	Article.Img = "/static/img/" + filename
	o.Insert(&Article)
	this.Redirect("/article/index", 302)

}

/**************************文章类型展示页面*******************************/
//文章类型展示界面
func (this *ArticleController) ShowAddType() {
	o := orm.NewOrm()
	var TypeNames []models.ArticleType
	o.QueryTable("ArticleType").All(&TypeNames)
	this.Data["TypeNames"] = TypeNames

	userName := this.GetSession("UserName")
	this.Data["userName"] = userName
	this.Layout = "articlelayout.html"
	this.TplName = "addType.html"

}

//处理添加类型
func (this *ArticleController) HandleAddType() {
	TypeName := this.GetString("typeName")
	if TypeName == "" {
		beego.Info("添加分类为空")
		this.Redirect("/article/addtype", 302)
		return
	}
	o := orm.NewOrm()
	var Typename models.ArticleType
	Typename.TypeName = TypeName
	err := o.Read(&Typename, "TypeName")
	if err != nil {
		o.Insert(&Typename)
	}
	this.Redirect("/article/addtype", 302)

}

func (this *ArticleController) DeleteAddType() {
	ArticleTypeName := this.GetString("typename")
	o := orm.NewOrm()
	var Articletype models.ArticleType
	Articletype.TypeName = ArticleTypeName
	err := o.Read(&Articletype, "TypeName")
	if err != nil {
		beego.Info("查询数据出问题")
		return
	}
	_, err = o.Delete(&Articletype)
	if err != nil {
		beego.Info("删除数据出问题")
		return
	}
	this.Redirect("/article/addtype", 302)

}
