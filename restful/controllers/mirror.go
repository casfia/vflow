package controllers

import (
	"github.com/astaxie/beego"
	"../../mirror"
	"fmt"
)


// Operations about object
type MirrorController struct {
	beego.Controller
	MirrorService mirror.Netflowv9Mirror
}


// @Title Get
// @Description find object by objectid
// @Param	objectId		path 	string	true		"the objectid you want to get"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [get]
func (o *MirrorController) Get() {
	fmt.Printf("call get method of mirror!\r\n")
	configs := mirror.Netflowv9Instance.GetConfig()
	sourceId := o.GetString("sourceId")

	//sourceId := o.Ctx.Input.Param(":sourceId")
	if sourceId != "" {
		for _,e := range configs {
			if e.Source == sourceId {
				o.Data["json"] = e
				o.ServeJSON()
				return
			}
		}
	}else{
		configs := mirror.Netflowv9Instance.GetConfig()

		fmt.Printf("serve all configs\r\n")
		o.Data["json"] = configs
		o.ServeJSON()
		return
	}
	o.Data["json"] = map[string]interface{}{}
	o.ServeJSON()
	return
}
