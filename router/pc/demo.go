package pc

import (
	"douyu/controller/pc/demo"
	"github.com/gin-gonic/gin"
)

func DemoRouter(api *gin.RouterGroup) {
	r := api.Group("demo")
	{
		// oss 回调
		r.GET("demo", demo.Demo)
	}
}
