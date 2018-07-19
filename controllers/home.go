package controllers

import (
	"github.com/astaxie/beego"
	"frontend/db"
	"time"
	"strconv"
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get()  {
	openid := this.GetString(":openid")
	l , err := db.Dia.Query("select register_reward from elastos_register_details where openid = '"+ openid +"'")
	if err != nil || l.Len() == 0{
		beego.Error(" error ", err , " or can not find openid ")
		this.TplName = "error.html"
		return
	}
	m := l.Front().Value.(map[string]string)
	if m["register_reward"] == "NULL" || m["register_reward"] == "" {
		this.Redirect("/registerUser/"+openid+"?"+strconv.Itoa(int(time.Now().Unix())),302)
		return
	}
	this.Data["reward"] = m["register_reward"]
	this.Data["openid"] = openid
	this.Data["elaWallet"] = beego.AppConfig.String("walletAddr")
	this.TplName = "home.tpl"
}
