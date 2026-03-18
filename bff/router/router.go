package router

import (
	"order/bff/handler/notify"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/callback", notify.Notify)
	r.POST("/notify/pay", notify.Notify)
	return r
}
