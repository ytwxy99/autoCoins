package policy

import (
	"fmt"

	"github.com/ytwxy99/autoCoins/interfaces"
	"github.com/ytwxy99/autoCoins/utils"
)

type Umbrella struct{}

func (*Umbrella) Target(args ...interface{}) interface{} {
	//db := args[0].(*gorm.DB)
	//sysConf := args[1].(*configuration.SystemConf)

	//futures, err := utils.ReadLines(sysConf.UmbrellaCsv)
	//if err != nil {
	//	logrus.Error("read lines error: %v", err)
	//}

	spotsK := (&interfaces.MarketArgs{
		CurrencyPair: "BTC_USDT",
		Interval:     -900,
		Level:        utils.Level1Day,
	}).SpotMarket()

	futuresK := (&interfaces.MarketArgs{
		CurrencyPair: "BTC_USDT",
		Interval:     -900,
		Level:        utils.Level1Day,
	}).FutureMarket()

	offset := offset(len(spotsK), len(futuresK))
	if offset != 0 {
		if len(spotsK) > len(futuresK) {
			for i, future := range futuresK {
				spotIndex := i + offset
				fmt.Println(spotsK[spotIndex][0], spotsK[spotIndex][2])
				fmt.Println(utils.GetData(int64(future.T)), future.C)
				fmt.Println((utils.StringToFloat64(spotsK[spotIndex][2]) - utils.StringToFloat64(future.C)) / utils.StringToFloat64(future.C))
			}
		} else {

		}
	} else {
	}

	return nil
}

func offset(a int, b int) int {
	if a < b {
		return b - a
	} else {
		return a - b
	}
}
