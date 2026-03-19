package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);comment:'用户名'"`
	Mobile   string `gorm:"type:char(11);unique;comment:'手机号'"`
	Password string `gorm:"type:varchar(32);comment:'密码'"`
}
