package inits

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gospacex-tz/srv/basic/config"
	"gospacex-tz/srv/handler/model"
)

func MysqlInit() {
	var err error
	data := config.GlobalConfig.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		data.User,
		data.Password,
		data.Host,
		data.Port,
		data.Database)
	config.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}
	fmt.Println("数据库连接成功")

	err = config.DB.AutoMigrate(&model.User{}, &model.Product{}, &model.Order{}, &model.OrderItem{}, &model.Inventory{}, &model.Category{})
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}
	fmt.Println("数据库迁移成功")
}
