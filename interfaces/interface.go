package interfaces

import (
	"strconv"

	c "github.com/ytwxy99/autoCoins/client"
	"github.com/ytwxy99/autoCoins/utils"
)

// get k market data
func Market(currencyPair string, beforeInterval int, interval string) [][]string {
	from := utils.GetOldTimeStamp(0, 0, beforeInterval)
	to := utils.GetNowTimeStamp()
	values := c.GetSpotCandlesticks(currencyPair, from, to, interval)

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
