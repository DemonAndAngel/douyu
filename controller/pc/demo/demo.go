package demo

import (
	"douyu/utils/response"
	"github.com/gin-gonic/gin"
)

func Demo(c *gin.Context) {
	v := map[string]interface{}{
		"demo": "demo",
	}
	response.MakeSuccessRespJson(c, v)
}