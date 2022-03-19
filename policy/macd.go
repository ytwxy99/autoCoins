package policy

import (
	"github.com/ytwxy99/autoCoins/interfaces"
	"github.com/ytwxy99/autoCoins/utils"
)

type MacdPolicy struct{}

// get macd data
func GetMacd(values [][]string, short int, long int, M int) []map[string]string {
	var ks []map[string]string
	for _, value := range values {
		m := make(map[string]string)
		m["t"] = value[0] // time
		m["v"] = value[1] // volume
		m["c"] = value[2] // close price
		m["h"] = value[3] // high price
		m["l"] = value[4] // low price
		m["o"] = value[5] // open price
		ks = append(ks, m)
	}

	emas := calcEma(ks, short)
	emaq := calcEma(ks, long)

	for index, k := range ks {
		k["diff"] = utils.Float32ToString(emas[index] - emaq[index])
	}

	for index, k := range ks {
		if index == 0 {
			k["dea"] = k["diff"]
		} else {
			dea := utils.StringToFloat32(ks[index-1]["dea"])
			diff := utils.StringToFloat32(k["diff"])
			deaIndex := (float32(M-1)*dea + 2*diff) / float32(M+1)
			k["dea"] = utils.Float32ToString(deaIndex)
		}
	}

	for _, k := range ks {
		diff := utils.StringToFloat32(k["diff"])
		dea := utils.StringToFloat32(k["dea"])
		k["macd"] = utils.Float32ToString(2 * (diff - dea))
	}

	return ks
}

// cal ema
func calcEma(values []map[string]string, num int) []float32 {
	var emaAll []float32

	for index, value := range values {
		if index == 0 {
			c := utils.StringToFloat32(value["c"])
			emaAll = append(emaAll, c)
			value["ema"] = utils.Float32ToString(c)
		} else {
			c := utils.StringToFloat32(value["c"])
			ema := utils.StringToFloat32(values[index-1]["ema"])
			value["ema"] = utils.Float32ToString(ema)

			ema = (float32(num-1) * ema) + (2*c)/float32(num+1)
			emaAll = append(emaAll, ema)
		}
	}

	return emaAll
}

// find macd buy point
func (*MacdPolicy) Target(args ...interface{}) interface{} {
	// convert specified type
	coin := args[0].(string)

	k4hValues := interfaces.Market(coin, -100, "4h")
	k15mValues := interfaces.Market(coin, -1, "15m")
	if k4hValues != nil && k15mValues != nil {
		k4hMacds := GetMacd(k4hValues, 12, 26, 9)
		k15mMacds := GetMacd(k15mValues, 12, 26, 9)
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
