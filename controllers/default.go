package controllers

import (
	"github.com/astaxie/beego"
	"github.com/skip2/go-qrcode"
	"github.com/astaxie/beego/httplib"
	"encoding/json"
)

type MainController struct {
	beego.Controller
}


func (c *MainController) Get() {
	req := httplib.Get(beego.AppConfig.String("GET_VALIDATION_CODE_URL"))
	retStr , err := req.String()
	if err != nil {
		beego.Error("request validation code error")
		c.TplName = "error.html"
		return
	}
	retMap :=make(map[string]interface{})
	err = json.Unmarshal([]byte(retStr),&retMap)
	if err != nil {
		beego.Error("request validation return result is not a json")
		c.TplName = "error.html"
		return
	}
	code := retMap["Result"].(map[string]interface{})["code"].(string)
	err = qrcode.WriteFile("https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx6bb613249b29f815&redirect_uri=https%3A%2F%2Ffun-test.elastos.org%2Flottery%2Fwechat%2F"+code+"&response_type=code&scope=snsapi_userinfo&state=3d6ae0a4035d839573b04816624a415e#wechat_redirect", qrcode.Medium, 256, "static/images/qr.png")
	if err != nil {
		beego.Error("error generating qrcode")
		c.TplName = "error.html"
		return
	}
	c.Data["Website"] = "elastos.org"
	c.Data["Email"] = "support@elastos.org"
	c.TplName = "welcome.tpl"
}

