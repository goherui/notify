package notify

import (
	"fmt"
	"net/http"
	"order/bff/basic/config"
	payment "order/proto/payment"

	"github.com/gin-gonic/gin"
)

func Notify(c *gin.Context) {
	orderSn := c.Query("out_trade_no")
	fmt.Println(orderSn)
	if orderSn == "" {
		orderSn = c.Query("orderSn")
	}
	if orderSn == "" {
		orderSn = c.PostForm("out_trade_no")
	}
	if orderSn == "" {
		orderSn = c.PostForm("orderSn")
	}
	if orderSn == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少订单号",
		})
		return
	}

	r, err := config.PaymentClient.HandlePaymentNotify(c, &payment.PaymentNotifyRequest{
		OrderSn: orderSn,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "处理支付回调失败",
		})
		return
	}
	if !r.Success {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "处理支付回调失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
	})
}
