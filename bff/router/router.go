package router

import (
	"github.com/gin-gonic/gin"
	"gospacex-tz/bff/handler/service"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.POST("product/add", service.ProductAdd)

	r.POST("order/add", service.OrderCreate)

	r.POST("notify/pay", service.NotifyPay)

	return r
}
