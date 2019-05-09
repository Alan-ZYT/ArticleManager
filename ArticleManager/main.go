package main

import (
	_ "ArticleManager/routers"
	"github.com/astaxie/beego"
	_ "ArticleManager/models"
)

func main() {
	beego.AddFuncMap("KeyAdd",KeyAdd)
	beego.AddFuncMap("Up",LastPage)
	beego.AddFuncMap("Down",NextPage)
	beego.Run()
}

func KeyAdd(key int) int {
	return key+1

}
func LastPage(index int)int{
	if index==1{
		return 1
	}
	return index-1
}
func NextPage(index int,num float64)int{
	if index==int(num){
		return int(num)
	}
	return index+1
}
