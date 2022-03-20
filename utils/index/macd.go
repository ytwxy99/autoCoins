package index

import "github.com/ytwxy99/autoCoins/utils"

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
