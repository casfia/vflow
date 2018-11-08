// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import(
	"github.com/astaxie/beego"
	"../controllers"
	"fmt"
	"../../vflow"
)

func Init(netflowv9 *vflow.NetflowV9) {
	fmt.Printf("router ininted.\n")
	beego.Router("/user", &controllers.UserController{})
	beego.Router("/policy", &controllers.PolicyController{})
	beego.Router("/rule", &controllers.RuleController{})

	controller := &controllers.StatsController{}
	controller.Init(*netflowv9)
	beego.Router("/stats",controller)
}
