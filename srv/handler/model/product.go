package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string  `gorm:"type:varchar(50);index;comment:'商品名称'"`
	Price       float64 `gorm:"type:decimal(10,2);comment:'商品价格'"`
	Images      string  `gorm:"type:varchar(100);comment:'商品图片'"`
	Description string  `gorm:"type:varchar(50);comment:'商品描述'"`
	CategoryId  int     `gorm:"type:int;index;comment:'分类id'"`
	Status      int     `gorm:"type:int;default:1;comment:'商品状态'"` /// 1上架  2下架
	Stock       int     `gorm:"type:int;comment:'库存'"`
}

func (p *Product) ProductAdd(db *gorm.DB) error {
	return db.Debug().Create(&p).Error
}

type Order struct {
	gorm.Model
	OrderSn   string `gorm:"type:varchar(30);comment:'订单号'"`
	UserId    int    `gorm:"type:bigint unsigned;index:idx_user_created,sort:asc;comment:'用户ID'"`
	PayType   int    `gorm:"type:int;comment:'支付方式:1-支付宝'"`
	AddressId int    `gorm:"type:int;comment:'收货地址ID'"`
	Status    int    `gorm:"type:int;comment:'订单状态:0-待支付,1-已支付,2-已完成,3-已超时,4-已取消'"`
}

type OrderItem struct {
	gorm.Model
	OrderId      int     `gorm:"type:int;index;comment:'订单ID'"`
	ProductId    int     `gorm:"type:int;comment:'商品ID'"`
	ProductName  string  `gorm:"type:varchar(50);comment:'商品名称'"`
	ProductPrice float64 `gorm:"type:decimal(10,2);comment:'商品单价'"`
	ProductImg   string  `gorm:"type:varchar(50);comment:'商品图片'"`
	ProductBio   string  `gorm:"type:varchar(50);comment:'商品详细'"`
	Quantity     int     `gorm:"type:int(11);comment:'购买数量'"`
}

type Inventory struct {
	gorm.Model
	ProductId int `gorm:"type:int;comment:'商品id'"`
	Stock     int `gorm:"type:int unsigned;default:0;comment:'库存数量'"`
}

type Category struct {
	gorm.Model
	Name string `gorm:"type:varchar(50);comment:'分类名称'"`
}
