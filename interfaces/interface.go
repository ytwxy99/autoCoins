package interfaces

import (
	"strconv"

	c "github.com/ytwxy99/autoCoins/client"
	"github.com/ytwxy99/autoCoins/utils"
)

type MarketArgs struct {
	CurrencyPair string
	Interval     int
	Level        string
}

// get market data
func (marketArgs *MarketArgs) Market() [][]string {
	from := utils.GetOldTimeStamp(0, 0, marketArgs.Interval)
	to := utils.GetNowTimeStamp()
	values := c.GetSpotCandlesticks(marketArgs.CurrencyPair, from, to, marketArgs.Level)

	if values != nil {
		for _, v := range values {
			timestamp, _ := strconv.ParseInt(v[0], 10, 64)
			v[0] = utils.GetData(timestamp)
		}

		return values
	} else {
		return nil
	}
}
