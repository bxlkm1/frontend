package routers

import (
	"frontend/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/qrlogin/:vldCode",&controllers.QrLoginController{})
    beego.Router("/index",&controllers.IndexController{})
	beego.Router("/submitAddr/:addr",&controllers.IndexController{},"get:SubmitAddr")
    beego.Router("/home/:openid",&controllers.HomeController{})
	beego.Router("/registerUser/:openid",&controllers.IndexController{},"get:RegisterUser")

}
