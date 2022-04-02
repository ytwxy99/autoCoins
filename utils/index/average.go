package index

import (
	"math"

	"github.com/ytwxy99/autoCoins/interfaces"
	"github.com/ytwxy99/autoCoins/utils"
)

type Average struct {
	CurrencyPair string
	Intervel     int
	Level        string
}

/**
 * @Description: five's average market index
 * @param: the level of markets which support:
 * 10s, 1m, 5m, 15m, 30m, 1h, 4h, 8h, 1d, 7d
 */
func (average *Average) FiveAverage(backOff bool) float64 {
	if !backOff {
		var sum float64
		markets := (&interfaces.MarketArgs{
			CurrencyPair: average.CurrencyPair,
			Interval:     average.Intervel,
			Level:        average.Level,
		}).SpotMarket()

		for _, market := range markets {
			sum += utils.StringToFloat64(market[2])
		}

		return sum / (math.Abs(float64(average.Intervel)) + 1)
	} else {
		var sum float64
		markets := (&interfaces.MarketArgs{
			CurrencyPair: average.CurrencyPair,
			Interval:     average.Intervel - 1,
			Level:        average.Level,
		}).SpotMarket()

		for index, market := range markets {
			if index == -(average.Intervel - 1) {
				continue
			}
			sum += utils.StringToFloat64(market[2])
		}

		return sum / (math.Abs(float64(average.Intervel)) + 1)
	}
}
