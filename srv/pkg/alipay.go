package pkg

import (
	"github.com/smartwalle/alipay/v3"
	"gospacex-tz/srv/basic/config"
	"strconv"
)

func AliPay(orderSn string, total float64) string {
	pay := config.GlobalConfig.AliPay
	privateKey := pay.PrivateKey           // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	alipayPublicKey := pay.AlipayPublicKey // 支付宝公钥
	appId := pay.AppId
	client, err := alipay.New(appId, privateKey, false)
	if err != nil {
		return ""
	}

	// 加载支付宝公钥
	if alipayPublicKey != "" {
		err = client.LoadAliPayPublicKey(alipayPublicKey)
		if err != nil {
			return ""
		}
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = pay.NotifyURL
	p.ReturnURL = pay.ReturnURL
	p.Subject = "京东在线支付"
	p.OutTradeNo = orderSn
	p.TotalAmount = strconv.FormatFloat(total, 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		return ""
	}

	// 这个 payURL 即是用于打开支付宝支付页面的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	var payURL = url.String()
	return payURL
}
