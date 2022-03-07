package interfaces

import (
	"github.com/gateio/gateapi-go/v6"
	c "github.com/ytwxy99/autoCoins/client"
	"github.com/ytwxy99/autoCoins/utils"
	"strconv"
)

// get k market data
func K(client *gateapi.APIClient, currencyPair string, beforeInterval int, interval string) [][]string {
	from := utils.GetOldTimeStamp(0, 0, beforeInterval)
	to := utils.GetNowTimeStamp()
	values := c.GetSpotCandlesticks(client, currencyPair, from, to, interval)

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
