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
)

func init() {
	fmt.Printf("router ininted.\n")
	beego.Router("/user", &controllers.UserController{})
	beego.Router("/mirror", &controllers.MirrorController{})
}
