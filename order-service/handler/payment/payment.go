package payment

import (
	"context"
	"errors"
	"fmt"
	"log"
	"order/order-service/basic/config"
	"order/order-service/model"
	"order/pkg"
	payment "order/proto/payment"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/gorm"
)

// server is used to implement helloworld.GreeterServer.
type Server struct {
	payment.UnimplementedPaymentServiceServer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) Alipay(_ context.Context, in *payment.AlipayRequest) (*payment.AlipayResponse, error) {
	orderSn := pkg.OrderSn()
	total := 0.0
	var items []*model.OrderItem
	for _, g := range in.Goods {
		var goods model.Goods
		err := goods.FindGoodsId(config.DB, g.GoodsID)
		if err != nil {
			return nil, errors.New("商品不存在")
		}
		suTotal := goods.Price * float64(g.Quantity)
		total += suTotal
		items = append(items, &model.OrderItem{
			GoodsID:  int(g.GoodsID),
			Title:    goods.GoodsName,
			Price:    goods.Price,
			Quantity: int(g.Quantity),
			SubTotal: suTotal,
		})
	}
	order := model.PaymentOrder{
		UserID:  int(in.UserId),
		OrderSn: orderSn,
		Total:   total,
		Status:  "待支付",
	}
	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return errors.New("创建订单失败")
		}
		for i := range items {
			items[i].PaymentOrderID = int(order.ID)
			if err := tx.Create(items[i]).Error; err != nil {
				return errors.New("创建订单商品失败")
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	Url := pkg.Alipay(orderSn, total)
	return &payment.AlipayResponse{
		Success: true,
		OrderSn: orderSn,
		Total:   float32(total),
		PayUrl:  Url,
	}, nil
}
func (s *Server) HandlePaymentNotify(_ context.Context, in *payment.PaymentNotifyRequest) (*payment.PaymentNotifyResponse, error) {
	if in.OrderSn == "" {
		return &payment.PaymentNotifyResponse{
			Success: false,
			Message: "订单号不能为空",
		}, nil
	}
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var order model.PaymentOrder
		err := order.FindOrderSn(tx, in.OrderSn)
		if err != nil {
			return errors.New("订单不存在")
		}
		if order.Status != "待支付" {
			return nil
		}
		var items []*model.OrderItem
		if err := tx.Where("payment_order_id=?", order.ID).Find(&items).Error; err != nil {
			return errors.New("查询订单商品失败")
		}
		for _, item := range items {
			var goods model.Goods
			err = goods.FindGoodsId(tx, int64(item.GoodsID))
			if err != nil {
				return errors.New("商品不存在")
			}
			if goods.Stock <= item.Quantity {
				return errors.New("库存不足")
			}
			if err = tx.Model(&goods).Update("stock", goods.Stock-item.Quantity).Error; err != nil {
				return errors.New("扣减库存失败")
			}
			if err := tx.Model(&order).Update("status", "已支付").Error; err != nil {
				return errors.New("更新订单状态失败")
			}
		}
		s.sendStockMsg(tx, order, items)
		return nil
	})
	if err != nil {
		return &payment.PaymentNotifyResponse{Success: false, Message: err.Error()}, nil
	}
	return &payment.PaymentNotifyResponse{Success: true, Message: "支付成功"}, nil
}
func (s *Server) sendStockMsg(tx *gorm.DB, order model.PaymentOrder, items []*model.OrderItem) {
	for _, item := range items {
		var goods model.Goods
		if tx.First(&goods, item.GoodsID).Error != nil {
			continue
		}
		msg := fmt.Sprintf("订单:%s,商品%d(%s),扣减%d,剩余%d",
			order.OrderSn, item.GoodsID, item.Title, item.Quantity, goods.Stock)
		rmq := pkg.NewRabbitMQSimple("stock_deduct_log")
		rmq.PublishSimple(msg)
		rmq.Destory()
	}
}

func StartStockConsumer() {
	logFile, _ := os.OpenFile("stock_deduct.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer logFile.Close()
	logger := log.New(logFile, "[库存扣减]", log.LstdFlags)
	rmq := pkg.NewRabbitMQSimple("stock_deduct_log")
	defer rmq.Destory()
	log.Println("启动库存扣减消费者...")
	rmq.ConsumeSimpleCallback(func(body []byte) {
		logMsg := fmt.Sprintf("%s - %s", time.Now().Format("2006-01-02 15:04:05"), string(body))
		logger.Println(logMsg)
		log.Printf("处理库存消息: %s", body)
	})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("停止库存扣减消费者...")
}
