package model

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name   string  `gorm:"column:name;not null;type:varchar(255);comment:商品名称" json:"name"`
	Price  float64 `gorm:"column:price;not null;type:decimal(10,2);comment:商品价格（使用decimal避免float精度问题）" json:"price"`
	Stock  int64   `gorm:"column:stock;not null;default:0;comment:商品库存" json:"stock"`
	Points int64   `gorm:"column:points;not null;default:0;comment:商品对应积分" json:"points"`
	Status int32   `gorm:"column:status;not null;type:tinyint;comment:商品状态（0-下架，1-上架）" json:"status"`
}
