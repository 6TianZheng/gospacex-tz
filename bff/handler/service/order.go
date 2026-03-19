package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gospacex-tz/srv/basic/config"
	"gospacex-tz/srv/handler/model"
	"net/http"
)

func NotifyPay(c *gin.Context) {
	c.Request.ParseForm()
	form := c.Request.PostForm
	fmt.Println("aaa", form)
	tradeStatus := c.Request.PostForm.Get("trade_status")
	outTradeNo := c.Request.PostForm.Get("out_trade_no")

	fmt.Println("tradeStatus:", tradeStatus, "outTradeNo:", outTradeNo)

	if tradeStatus != "TRADE_SUCCESS" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "支付失败",
			"code": 400,
		})
		return
	}

	if outTradeNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "订单不存在",
			"code": 400,
		})
		return
	}

	var order model.Order
	err := config.DB.Where("order_sn = ?", outTradeNo).First(&order).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "订单不存在",
			"code": 400,
		})
		return
	}

	if order.Status == 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "订单已处理",
			"code": 400,
		})
		return
	}

	tx := config.DB.Begin()

	order.Status = 2
	err = tx.Save(&order).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "订单状态更新失败",
			"code": 400,
		})
		return
	}

	var orderItems []model.OrderItem
	err = tx.Where("order_id = ?", order.ID).Find(&orderItems).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "订单不存在",
			"code": 400,
		})
		return
	}

	for _, item := range orderItems {
		var product model.Product
		err = tx.First(&product, item.ProductId).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":  "商品不存在",
				"code": 400,
			})
			return
		}
		product.Stock -= item.Quantity
		err = tx.Save(&product).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":  "库存扣减失败",
				"code": 400,
			})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"msg":  "库存扣减成功",
		"code": 200,
	})
	return
}
