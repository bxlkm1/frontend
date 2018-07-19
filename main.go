package main

import (
	_ "frontend/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.SetStaticPath("/static","static")
	beego.Run()
}

