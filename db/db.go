package db

import (
	"github.com/xiaomingfuckeasylife/job/db"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

var Dia *db.Dialect
var Orm orm.Ormer

func init()  {
	Dia = &db.Dialect{}
	Dia.Create(beego.AppConfig.String("dbDriverName"),beego.AppConfig.String("dbDriverSource"))
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.Debug = true
	orm.RegisterDataBase("frontend", beego.AppConfig.String("dbDriverName"), beego.AppConfig.String("dbDriverSource"))
}
