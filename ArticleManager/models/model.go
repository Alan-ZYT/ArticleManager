package models

import(
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"time"
)

type User struct{
	Id int
	UserName string `orm:"size(30)"`
	PassWord string `orm:"size(40)"`
	Articles []*Article `orm:"rel(m2m)"`
}

type Article struct {
	Id int	`orm:"pk;auto"`
	Title string 	`orm:"size(40);unique"`
	Content string  `orm:"size(2000)"`
	Time time.Time	`orm:"type(datetime);auto_now_add"`
	ReadCount int	`orm:"default(0)"`
	Img string		`orm:"null"`

	ArticleType *ArticleType `orm:"rel(fk);null;on_delete(set_null)"`
	Users []*User `orm:"reverse(many)"`
}

type ArticleType struct{
	Id int
	TypeName string `orm:"unique;size(200)"`
	Articles []*Article `orm:"reverse(many)"`

}

func init(){
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/test1")
	orm.RegisterModel(new(User),new(ArticleType),new(Article))
	orm.RunSyncdb("default",false,true)
}
