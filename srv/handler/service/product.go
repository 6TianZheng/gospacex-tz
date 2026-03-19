package service

import (
	"context"
	"errors"
	"gospacex-tz/srv/basic/config"
	__ "gospacex-tz/srv/basic/proto"
	"gospacex-tz/srv/handler/model"
)

// server is used to implement helloworld.GreeterServer.
type Server struct {
	__.UnimplementedProductServer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) ProductAdd(_ context.Context, in *__.ProductAddReq) (*__.ProductAddResp, error) {

	product := model.Product{
		Name:        in.Name,
		Price:       float64(in.Price),
		Images:      in.Images,
		Description: in.Description,
		CategoryId:  int(in.CategoryId),
	}

	err := product.ProductAdd(config.DB)
	if err != nil {
		return nil, errors.New("商品上架失败")
	}

	return &__.ProductAddResp{
		Msg:  "商品上架成功",
		Code: 200,
	}, nil
}
