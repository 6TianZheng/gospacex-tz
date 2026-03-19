package service

import (
	"context"
	"errors"
	"gospacex-tz/srv/basic/config"
	__ "gospacex-tz/srv/basic/proto"
	"gospacex-tz/srv/handler/model"
	"gospacex-tz/srv/pkg"
)

func (s *Server) OrderItemAdd(_ context.Context, in *__.OrderItemAddReq) (*__.OrderItemAddResp, error) {
	orderSn := pkg.OrderSn()
	total := 0.0

	var list []*model.OrderItem

	for _, item := range in.List {
		var product model.Product
		err := config.DB.Where("id = ?", item.ProductId).First(&product).Error
		if err != nil {
			return nil, errors.New("商品不存在")
		}

		if product.Status != 1 {
			return nil, errors.New("商品未上架")
		}

		if item.Quantity > int64(product.Stock) {
			return nil, errors.New("库存不足")
		}

		total += product.Price * float64(item.Quantity)

		list = append(list, &model.OrderItem{
			ProductId:    int(int64(product.ID)),
			ProductName:  product.Name,
			ProductPrice: product.Price,
			ProductImg:   product.Images,
			ProductBio:   product.Description,
			Quantity:     int(item.Quantity),
		})
	}

	order := model.Order{
		OrderSn:   orderSn,
		UserId:    int(in.UserId),
		PayType:   int(in.PayType),
		AddressId: int(in.AddressId),
		Status:    1,
	}

	err := config.DB.Create(&order).Error
	if err != nil {
		return nil, errors.New("订单创建失败")
	}

	for i := range list {
		list[i].OrderId = int(int64(order.ID))
	}

	err = config.DB.Create(&list).Error
	if err != nil {
		return nil, errors.New("订单明细创建失败")
	}

	payUrl := pkg.AliPay(orderSn, total)

	return &__.OrderItemAddResp{
		OrderSn: orderSn,
		PayUrl:  payUrl,
		Total:   float32(total),
	}, nil
}
