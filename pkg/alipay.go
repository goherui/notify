package pkg

import (
	"fmt"
	"order/order-service/basic/config"
	"strconv"

	"github.com/smartwalle/alipay/v3"
)

func Alipay(orderSn string, total float64) string {
	var privateKey = config.GlobalConfig.AliPay.Key // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	appId := config.GlobalConfig.AliPay.AppId
	client, err := alipay.New(appId, privateKey, false)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var p = alipay.TradePagePay{}
	p.NotifyURL = config.GlobalConfig.AliPay.NotifyPay
	p.ReturnURL = config.GlobalConfig.AliPay.Return
	p.Subject = "支付宝收起"
	p.OutTradeNo = orderSn
	p.TotalAmount = strconv.FormatFloat(total, 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		fmt.Println(err)
	}

	// 这个 payURL 即是用于打开支付宝支付页面的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	var payURL = url.String()
	fmt.Println(payURL)
	return payURL
}
