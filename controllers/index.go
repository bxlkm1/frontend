package controllers

import (
	"github.com/astaxie/beego"
	"frontend/db"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"encoding/hex"
	"strconv"
	"github.com/astaxie/beego/httplib"
	"math/rand"
	"math"
	"fmt"
	"github.com/elastos/Elastos.ELA.Utility/common"
	"crypto/sha256"
)

type IndexController struct {
	beego.Controller
}

type RetMsg struct {
	Result interface{}
	Error  int
	Desc   string
}
// @router / [get]
func (c *IndexController) Get()  {
	openid := c.GetString("openid")
	isModify , err := c.GetInt("isModify")
	if err == nil && isModify == 1 {
		req := httplib.Get(beego.AppConfig.String("GET_VALIDATION_CODE_URL"))
		retStr , err := req.String()
		if err != nil {
			beego.Error("request validation code error")
			c.Data["msg"] = "request validation code error"
			c.TplName = "error.html"
			return
		}
		retMap :=make(map[string]interface{})
		err = json.Unmarshal([]byte(retStr),&retMap)
		if err != nil {
			beego.Error("request validation return result is not a json")
			c.Data["msg"] = "request validation return result is not a json"
			c.TplName = "error.html"
			return
		}
		code := retMap["Result"].(map[string]interface{})["code"].(string)
		c.Data["openid"] = openid
		c.Data["vldCode"] = code
		c.Data["registed"] = 0
		l , err := db.Dia.Query("select wallet_addr from elastos_members where openid = '" + openid +"'")
		if err != nil || l.Len() == 0 {
			beego.Error(err,l)
			c.Data["msg"] = "no such user"
			c.TplName = "error.html"
			return
		}
		c.Data["addr"] = l.Front().Value.(map[string]string)["wallet_addr"]
		c.Data["walletRegister"] = beego.AppConfig.String("walletRegister")
		c.TplName = "index.tpl"
		return
	}
	vldCode := c.GetString("code")
	if openid == "" || vldCode == ""{
		beego.Error("openid || validation Code can not be blank ",openid,vldCode)
		c.TplName = "error.html"
		return
	}
	l , err := db.Dia.Query("select * from elastos_register_details where openid = '"+ openid +"'")
	if err != nil {
		beego.Error(" error ", err)
		c.TplName = "error.html"
		return
	}

	if l.Len() > 0 {
		model := l.Front().Value.(map[string]string)
		c.Data["registed"] = 1;
		rewarded := model["status"]
		if rewarded == "1" {
			c.Redirect("./home/"+openid+"?"+strconv.Itoa(int(time.Now().Unix())),302)
			return
		}
	}else {
		c.Data["walletRegister"] = beego.AppConfig.String("walletRegister")
		c.Data["registed"] = 0
	}
	c.Data["openid"] = openid
	c.Data["vldCode"] = vldCode
	c.TplName = "index.tpl"
}


func (this *IndexController) SubmitAddr() {

	addr := this.GetString(":addr")
	vldCode := this.GetString("vldCode")
	openid  := this.GetString("openid")
	if addr == "" || vldCode == "" || openid == ""{
		beego.Error("addr or vldCode or openid can not be blank",addr,vldCode,openid)
		this.TplName = "error.html"
		return
	}

	_ , err := common.Uint168FromAddressWithCheck(addr)
	if err != nil {
		this.Data["msg"] = "Invalid ELA Address"
		this.TplName = "error.html"
		return
	}

	l , err := db.Dia.Query("select * from elastos_members where openid = '" + openid +"'")
	if err != nil {
		beego.Error(" error ", err)
		this.Data["msg"] = err.Error()
		this.TplName = "error.html"
		return
	}
	if l.Len() == 0 {
		this.Data["msg"] = "can not find such user"
		this.TplName = "error.html"
		return
	}
	model := l.Front().Value.(map[string]string)
	isExist := true
	if model["wallet_addr"] == "NULL" || model["wallet_addr"] == "" {
		isExist = false
	}
	// only check validation code when it is the first time binding address
	if isExist == false {
		l , err =db.Dia.Query("select * from elastos_info where vldCode = '" +vldCode +"'")
		if err != nil {
			beego.Error(" error ", err)
			this.Data["msg"] = err.Error()
			this.TplName = "error.html"
			return
		}
	}
	if l.Len() == 0 {
		this.Data["msg"] = "validation code is outof data please Rescan the QR code"
		this.TplName = "error.html"
		return
	}


	if model["wallet_addr"] == "" || model["wallet_addr"] == "NULL" {
		_ , err = db.Dia.Exec("insert into elastos_register_details(openId) values('"+openid+"')")
		if err != nil {
			beego.Error(err.Error())
			this.Data["msg"] = err.Error()
			this.TplName = "error.html"
			return
		}

	}
	_ , err = db.Dia.Exec("update elastos_members set wallet_addr = '" + addr +"' where openid = '" + openid +"'")
	if err != nil {
		beego.Error(err.Error())
		this.Data["msg"] = err.Error()
		this.TplName = "error.html"
		return
	}
	if isExist {
		this.Redirect("/home/"+openid+"?"+strconv.Itoa(int(time.Now().Unix())),302)
		return
	}
	this.Data["registed"] = 1
	this.Data["openid"] = model["openid"]
	this.TplName = "index.tpl"
}

func (this *IndexController) RegisterUser(){
	openid := this.GetString(":openid")
	if openid == "" {
		beego.Error("openid can not be blank ",openid)
		this.Data["json"] = RetMsg{Error:-2}
		this.ServeJSON()
		return
	}
	l , err := db.Dia.Query("select * from elastos_register_details where openid = '" + openid +"'")
	if err != nil {
		beego.Error(" error ", err)
		this.Data["json"] = RetMsg{Desc:err.Error(),Error:-2}
		this.ServeJSON()
		return
	}
	if l.Len() == 0 {
		this.Data["json"] = RetMsg{Desc:"can not find such user",Error:-2}
		this.ServeJSON()
		return
	}
	model := l.Front().Value.(map[string]string)
	l , err = db.Dia.Query("select * from elastos_addresses where status = 1")
	if err != nil || l.Len() == 0 {
		beego.Error("initial address is not ready " , err, l.Len())
		this.Data["json"] = RetMsg{Desc:"initial address is not ready ",Error:-2}
		this.ServeJSON()
		return
	}
	sendAddrModel := l.Front().Value.(map[string]string)
	l , err = db.Dia.Query("select wallet_addr from elastos_members where openid = '"+openid +"'")
	if err != nil || l.Len() == 0{
		beego.Error("can not find suh member " , err, l.Len())
		this.Data["json"] = RetMsg{Desc:"can not find suh member ",Error:-2}
		this.ServeJSON()
		return
	}
	userwalletModel := l.Front().Value.(map[string]string)

	if model["register_info_tx"] == "" || model["register_info_tx"] == "NULL" {
		raw := userwalletModel["wallet_addr"]+"&"+openid+"&"+strconv.Itoa(int(time.Now().Unix()))
		s := sha256.Sum256([]byte(raw))
		memo := "chinajoy"+hex.EncodeToString(s[:])
		body := `{
			"Action":"transfer",
			"Version":"1.0.0",
			"Data":
				{"senderAddr":"` + sendAddrModel["publicAddr"] + `",
   				 "senderPrivateKey":"` + sendAddrModel["privKey"] + `",
				 "memo":"`+memo+`",
				 "receiver":[{"address":"` + beego.AppConfig.String("ReceivingPubAddr") + `","amount":"0.00000001"}]
				}
			}`
		beego.Info("request body " , body)
		r := strings.NewReader(body)
		rsp, err := http.Post(beego.AppConfig.String("SendTransfer"), "application/json", r)
		if err != nil {
			beego.Error(err)
			this.Data["json"] = RetMsg{Desc:err.Error(),Error:-2}
			this.ServeJSON()
			return
		}
		bytes, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			beego.Error(err)
			this.Data["json"] = RetMsg{Desc:err.Error(),Error:-2}
			this.ServeJSON()
			return
		}
		beego.Debug("ret Msg : %s \n", string(bytes))
		var ret map[string]interface{}
		err = json.Unmarshal(bytes, &ret)
		if err != nil {
			beego.Error(err)
			this.Data["json"] = RetMsg{Desc:err.Error(),Error:-2}
			this.ServeJSON()
			return
		}
		if ret["error"].(float64) != 0 {
			this.Data["json"] = RetMsg{Desc:string(bytes),Error:-2}
			this.ServeJSON()
			return
		}
		_ , err = db.Dia.Exec("update elastos_register_details set register_info_tx = '"+ ret["result"].(string) +"', register_memo_hash = '" + memo + "',register_memo_raw='" + raw+"' where openid = '" + openid +"'")
		if err != nil  {
			this.Data["json"] = RetMsg{Desc:err.Error(),Error:-2}
			this.ServeJSON()
			return
		}
		this.Data["json"] = RetMsg{Desc:"Registed info",Error:-1}
		this.ServeJSON()
	}else if model["register_reward_tx"] == "" || model["register_reward_tx"] == "NULL" {
		tx := model["register_info_tx"]
		l , err = db.Dia.Query("select * from (select * from elastos_txblock where height = ( select height from elastos_txblock where txid = '" + tx +"')) a left join elastos_register_details b on a.txid = b.register_info_tx left join elastos_members c on b.openId = c.openid")
		if err != nil || l.Len() == 0 {
			this.Data["json"] = RetMsg{Desc:"block has not been sync yet",Error:-1}
			this.ServeJSON()
			return
		}
		rand.Seed(time.Now().Unix())
		regReward := math.Round(float64((rand.Int31n(10)) + 1) * 0.01 * 100000000/float64(l.Len()))*0.00000001
		var receivInfo string
		var i int
		for e := l.Front() ; e != nil ; e = e.Next() {
			v := e.Value.(map[string]string)

			if i == 0 && l.Len() != 1 {
				receivInfo += `[{"address":"`+v["wallet_addr"]+`","amount":"`+fmt.Sprintf("%.8f",regReward)+`"},`
			}else if(i == 0 && l.Len() == 1){
				receivInfo += `[{"address":"`+v["wallet_addr"]+`","amount":"`+fmt.Sprintf("%.8f",regReward)+`"}]`
			}else if i != l.Len() - 1  {
				receivInfo += `{"address":"`+v["wallet_addr"]+`","amount":"`+fmt.Sprintf("%.8f",regReward)+`"},`
			}else {
				receivInfo += `{"address":"`+v["wallet_addr"]+`","amount":"`+fmt.Sprintf("%.8f",regReward)+`"}]`
			}
			i++
		}


		body := `{
			"Action":"transfer",
			"Version":"1.0.0",
			"Data":
				{"senderAddr":"` + beego.AppConfig.String("PlateformAddr") + `",
   				 "senderPrivateKey":"` + beego.AppConfig.String("PlateformPrivkey") + `",
				 "memo":"chinajoy registery reward",
				 "receiver":`+receivInfo+`
				}
			}`
		beego.Info("request body " , body)
		r := strings.NewReader(body)
		rsp, err := http.Post(beego.AppConfig.String("SendTransfer"), "application/json", r)
		if err != nil {
			beego.Error(err)
			this.Data["json"] = RetMsg{Desc:err.Error(),Error:-2}
			this.ServeJSON()
			return
		}
		bytes, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			beego.Error(err)
			this.Data["json"] = RetMsg{Desc:err.Error(),Error:-2}
			this.ServeJSON()
			return
		}
		beego.Debug("ret Msg : %s \n", string(bytes))
		var ret map[string]interface{}
		err = json.Unmarshal(bytes, &ret)
		if err != nil {
			beego.Error(err)
			this.Data["json"] = RetMsg{Desc:err.Error(),Error:-2}
			this.ServeJSON()
			return
		}
		if ret["error"].(float64) != 0 {
			this.Data["json"] = RetMsg{Desc:string(bytes),Error:-2}
			this.ServeJSON()
			return
		}
		for e := l.Front() ; e != nil ; e = e.Next() {
			v := e.Value.(map[string]string)
			_  ,err = db.Dia.Exec("update elastos_register_details set register_reward = " + fmt.Sprintf("%.8f",regReward) + " , register_reward_tx = '" + ret["result"].(string) +"' where openid = '" + v["openId"]+"'")
			if err != nil {
				beego.Error("Error updating elastos_register_details register reward and tx " , err)
			}
		}
		_ , err = db.Dia.Exec(" update elastos_addresses set status = 0 where status = 1")
		if err != nil {
			beego.Error(err)
		}
		id , err := db.Dia.Exec(" update elastos_addresses set status = 0 where id=" + sendAddrModel["id"]+"+1")
		if err != nil {
			beego.Error(err)
		}else if id == 0{
			_ , err := db.Dia.Exec(" update elastos_addresses set status = 0 where id = (select id from elastos_addresses limit 1)")
			if err != nil {
				beego.Error(err)
			}
		}
		if err != nil {
			beego.Error(err)
		}
		this.Data["json"] = RetMsg{Desc:string(bytes),Error:-1}
		this.ServeJSON()

	}else {
		if i,err := strconv.Atoi(model["status"]) ; err != nil && i == 1 {
			this.Data["json"] = RetMsg{Error:0}
			this.ServeJSON()
			return
		}

		tx := model["register_reward_tx"]
		req := httplib.Get(beego.AppConfig.String("GetTransactionByHash") +"/" + tx)
		ret := make(map[string]interface{})
		err := req.ToJSON(&ret)
		beego.Info("getTxByhash ret " , ret)
		if err != nil {
			this.Data["json"] = RetMsg{Desc:err.Error(),Error:-2}
			this.ServeJSON()
			return
		}
		if ret["Error"].(float64) == 0 {
			// update status of register details
			_ , err = db.Dia.Exec("update elastos_register_details set status = 1 where register_reward_tx = '" + tx +"'")
			if err != nil {
				beego.Error(err)
				this.Data["json"] = RetMsg{Desc:" update error ",Error:-1}
				this.ServeJSON()
				return
			}
			this.Data["json"] = RetMsg{Error:0}
			this.ServeJSON()
			return
		}
		this.Data["json"] = RetMsg{Error:-1}
		this.ServeJSON()
	}
}