package controllers

import (
	"github.com/astaxie/beego"
	"../../mirror"
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
	configs := mirror.Netflowv9Instance.GetConfig()
	sourceId := o.Ctx.Input.Param(":sourceId")
	if sourceId != "" {
		for _,e := range configs {
			if e.Source == sourceId {
				o.Data["configs"] = e
				o.ServeJSON()
				return
			}
		}
	}else{
		o.Data["configs"] = configs
		o.ServeJSON()
	}
	o.Data["configs"] = map[string]interface{}{}
	o.ServeJSON()
	return
}