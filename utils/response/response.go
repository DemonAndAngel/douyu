package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


func MakeFailRespJson(c *gin.Context, code int, msg string) {
	MakeResponseJson(c, code, msg, nil)
	return
}

func MakeSuccessRespJson(c *gin.Context, data interface{}) {
	MakeResponseJson(c, http.StatusOK, "success", data)
	return
}

func MakeResponseJson(c *gin.Context, code int, msg string, data interface{}) {
	MakeBaseResponseJson(c, gin.H{
		"meta": gin.H{
			"code": code,
			"msg":  msg,
		},
		"data": data,
	})
	return
}

func MakeBaseResponseJson(c *gin.Context, resp interface{}) {
	c.JSON(http.StatusOK, resp)
	c.Abort()
	return
}
