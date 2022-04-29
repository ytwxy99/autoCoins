package interfaces

import (
	"context"
	"strconv"

	"github.com/gateio/gateapi-go/v6"

	"github.com/ytwxy99/autoCoins/pkg/client"
	"github.com/ytwxy99/autoCoins/utils"
)

type MarketArgs struct {
	CurrencyPair string
	Interval     int
	Level        string
}

type Future struct {
	Settle string
}

func (marketArgs *MarketArgs) SpotMarket() [][]string {
	from := utils.GetOldTimeStamp(0, 0, marketArgs.Interval)
	to := utils.GetNowTimeStamp()
	values := client.GetSpotCandlesticks(marketArgs.CurrencyPair, from, to, marketArgs.Level)

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

func (marketArgs *MarketArgs) FutureMarket() []gateapi.FuturesCandlestick {
	from := utils.GetOldTimeStamp(0, 0, marketArgs.Interval)
	to := utils.GetNowTimeStamp()
	futures := client.GetFutureCandlesticks(marketArgs.CurrencyPair, from, to, marketArgs.Level)

	if futures != nil {
		return futures
	} else {
		return nil
	}
}

func (future *Future) GetAllFutures(ctx context.Context) ([]gateapi.Contract, error) {
	contracts, _, err := client.Client.FuturesApi.ListFuturesContracts(ctx, future.Settle)
	if err != nil {
		return nil, err
	}

	return contracts, err
}
