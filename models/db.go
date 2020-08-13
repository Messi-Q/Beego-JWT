package models

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.Debug = true
	if err := orm.RegisterDriver("mysql", orm.DRMySQL); err != nil {
		logs.Error(err.Error())
	}

	orm.RegisterModel(new(User))

	dbuser := beego.AppConfig.String("mysql::user")
	dbpass := beego.AppConfig.String("mysql::pass")
	dbhost := beego.AppConfig.String("mysql::host")
	dbport := beego.AppConfig.String("mysql::port")
	dbname := beego.AppConfig.String("mysql::db")

	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", dbuser, dbpass, dbhost, dbport, dbname)
	logs.Info("Will connect to mysql url", dbURL)

	if err := orm.RegisterDataBase("default", "mysql", dbURL); err != nil {
		logs.Error(err.Error())
		panic(err.Error())
	}

}
