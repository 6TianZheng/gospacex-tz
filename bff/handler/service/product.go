package service

import (
	"github.com/gin-gonic/gin"
	"gospacex-tz/bff/basic/config"
	"gospacex-tz/bff/handler/request"
	__ "gospacex-tz/srv/basic/proto"
	"net/http"
)

func ProductAdd(c *gin.Context) {
	var form request.ProductAdd
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "参数有误",
			"code": 400,
		})
		return
	}

	r, err := config.ProductClient.ProductAdd(c, &__.ProductAddReq{
		Name:        form.Name,
		Price:       float32(form.Price),
		Images:      form.Images,
		Description: form.Description,
		CategoryId:  int64(form.CategoryId),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "参数有误",
			"code": 400,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  r.Msg,
		"code": r.Code,
	})

}

func OrderCreate(c *gin.Context) {
	var req request.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构建订单商品列表
	var list []*__.OrderItem
	for _, item := range req.List {
		list = append(list, &__.OrderItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	resp, err := config.ProductClient.OrderItemAdd(c, &__.OrderItemAddReq{
		UserId:    req.UserId,
		PayType:   req.PayType,
		AddressId: req.AddressId,
		List:      list,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"OrderSn": resp.OrderSn,
		"Url":     resp.PayUrl,
		"Total":   resp.Total,
	})
	return
}
