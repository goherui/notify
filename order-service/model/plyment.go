package model

import "gorm.io/gorm"

type PaymentOrder struct {
	gorm.Model
	UserID  int     `gorm:"type:int;comment:用户id"`
	OrderSn string  `gorm:"type:char(32);comment:订单号"`
	Total   float64 `gorm:"type:decimal(10,2);comment:价格"`
	Status  string  `gorm:"type:varchar(30);comment:订单状态"`
}

func (o *PaymentOrder) FindOrderSn(tx *gorm.DB, sn string) error {
	return tx.Where("order_sn = ?", sn).First(o).Error
}

type OrderItem struct {
	gorm.Model
	PaymentOrderID int     `gorm:"type:int;comment:订单id"`
	GoodsID        int     `gorm:"type:int;comment:商品id"`
	Title          string  `gorm:"type:varchar(30);comment:商品名称"`
	Price          float64 `gorm:"type:decimal(10,2);comment:价格"`
	Quantity       int     `gorm:"type:int(11);comment:购买数量"`
	SubTotal       float64 `gorm:"type:decimal(10,2);comment:小计"`
}
type Goods struct {
	gorm.Model         // 包含ID、CreatedAt、UpdatedAt、DeletedAt字段
	GoodsName  string  `gorm:"type:varchar(100);not null;comment:商品名称"`
	Price      float64 `gorm:"type:decimal(10,2);not null;comment:商品单价"`
	Stock      int     `gorm:"type:int;default:0;comment:商品库存"`
	Status     int     `gorm:"type:tinyint;default:1;comment:商品状态 1-上架 2-下架 0-删除"`
}

func (g *Goods) FindGoodsId(db *gorm.DB, id int64) error {
	return db.Where("id = ?", id).First(&g).Error
}
