package routers

import (
	"Beego-Jwt/controllers"
	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/api",

		beego.NSNamespace("/user",
			beego.NSRouter("/login", &controllers.UserController{}, "post:Login"),
			beego.NSRouter("/register", &controllers.UserController{}, "post:CreateUser"),
		),
	)
	beego.AddNamespace(ns)
}
