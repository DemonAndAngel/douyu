package router

import (
	"douyu/router/pc"
	"github.com/gin-gonic/gin"
)

func InitRouter(app *gin.Engine) {
	api := app.Group("api")
	{
		// pc端接口
		pcApi := api.Group("pc")
		{
			pc.DemoRouter(pcApi)
		}
	}
}
