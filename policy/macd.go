package policy

import (
	"github.com/ytwxy99/autoCoins/interfaces"
	"github.com/ytwxy99/autoCoins/utils"
	"github.com/ytwxy99/autoCoins/utils/index"
)

type MacdPolicy struct{}

// find macd buy point
func (*MacdPolicy) Target(args ...interface{}) interface{} {
	// convert specified type
	coin := args[0].(string)

	k4hValues := interfaces.Market(coin, -100, "4h")
	k15mValues := interfaces.Market(coin, -1, "15m")
	if k4hValues != nil && k15mValues != nil {
		k4hMacds := index.GetMacd(k4hValues, 12, 26, 9)
		k15mMacds := index.GetMacd(k15mValues, 12, 26, 9)
		// fetch now local time macd data
		nowK4h := len(k4hMacds) - 1
		nowK15m := len(k15mMacds) - 1

		if nowK4h < 5 || nowK15m < 5 {
			return false
		}

		// set conditions
		H4Up := utils.StringToFloat32(k4hMacds[nowK4h]["c"]) > utils.StringToFloat32(k4hMacds[nowK4h]["o"])
		H4Up1 := utils.StringToFloat32(k4hMacds[nowK4h-1]["c"]) > utils.StringToFloat32(k4hMacds[nowK4h-2]["c"])*1.1
		H4Up2 := utils.StringToFloat32(k4hMacds[nowK4h-2]["c"]) > utils.StringToFloat32(k4hMacds[nowK4h-3]["c"])*1.1

		macdValue := utils.StringToFloat32(k4hMacds[nowK4h]["macd"])
		priceDiff := (utils.StringToFloat32(k4hMacds[nowK4h]["c"]) - utils.StringToFloat32(k4hMacds[nowK4h]["o"])) / utils.StringToFloat32(k4hMacds[nowK4h]["o"])
		// Increasing trading 15min macd
		k15MacDiff := utils.StringToFloat32(k15mMacds[nowK15m]["macd"]) >= utils.StringToFloat32(k15mMacds[nowK15m-1]["macd"])
		k15Macd := utils.StringToFloat32(k15mMacds[nowK15m]["macd"])

		return macdValue > 0 && H4Up && H4Up1 && H4Up2 && priceDiff >= 0.05 && k15MacDiff && k15Macd > 0
	}

	return false
}
