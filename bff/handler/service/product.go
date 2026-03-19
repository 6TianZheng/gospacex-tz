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
