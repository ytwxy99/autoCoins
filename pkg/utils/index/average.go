package index

import (
	"math"

	"github.com/ytwxy99/autocoins/pkg/interfaces"
	"github.com/ytwxy99/autocoins/pkg/utils"
)

type Average struct {
	CurrencyPair string
	Level        string
	MA           int
}

/**
 * @Description: five's average market index
 * @param: the level of markets which support:
 * 10s, 1m, 5m, 15m, 30m, 1h, 4h, 8h, 1d, 7d
 */
func (average *Average) Average(backOff bool) float64 {
	var sum float64
	if average.MA == utils.MA21 && average.Level == utils.Level4Hour {
		intervel := math.Ceil(float64(-4*average.MA/24)) - 1 //向上取整
		if !backOff {
			markets := (&interfaces.MarketArgs{
				CurrencyPair: average.CurrencyPair,
				Interval:     int(intervel),
				Level:        average.Level,
			}).SpotMarket()

			if markets == nil {
				return 0
			}

			for i, market := range markets {
				if i >= 4 {
					sum = sum + utils.StringToFloat64(market[2])
				}
			}

			return sum / float64(len(markets)-4)
		} else {
			markets := (&interfaces.MarketArgs{
				CurrencyPair: average.CurrencyPair,
				Interval:     int(intervel),
				Level:        average.Level,
			}).SpotMarket()

			if markets == nil {
				return 0
			}

			for i, market := range markets {
				if i == len(markets)-1 {
					continue
				}
				if i >= 3 {
					sum = sum + utils.StringToFloat64(market[2])
				}
			}

			return sum / float64(len(markets)-4)
		}
	} else if average.MA == utils.MA10 && average.Level == utils.Level4Hour {
		intervel := math.Ceil(float64(-4*average.MA/24)) - 1 //向上取整
		if !backOff {
			markets := (&interfaces.MarketArgs{
				CurrencyPair: average.CurrencyPair,
				Interval:     int(intervel),
				Level:        average.Level,
			}).SpotMarket()

			if markets == nil {
				return 0
			}

			for i, market := range markets {
				if i >= 3 {
					sum = sum + utils.StringToFloat64(market[2])
				}
			}

			return sum / float64(len(markets)-3)
		} else {
			markets := (&interfaces.MarketArgs{
				CurrencyPair: average.CurrencyPair,
				Interval:     int(intervel),
				Level:        average.Level,
			}).SpotMarket()

			if markets == nil {
				return 0
			}

			for i, market := range markets {
				if i == len(markets)-1 {
					continue
				}
				if i >= 2 {
					sum = sum + utils.StringToFloat64(market[2])
				}
			}

			return sum / float64(len(markets)-3)
		}
	} else if average.MA == utils.MA5 && average.Level == utils.Level4Hour {
		intervel := math.Ceil(float64(-4*average.MA/24)) - 1 //向上取整
		if !backOff {
			markets := (&interfaces.MarketArgs{
				CurrencyPair: average.CurrencyPair,
				Interval:     int(intervel),
				Level:        average.Level,
			}).SpotMarket()

			if markets == nil {
				return 0
			}

			for i, market := range markets {
				if i >= 2 {
					sum = sum + utils.StringToFloat64(market[2])
				}
			}

			return sum / float64(len(markets)-2)
		} else {
			markets := (&interfaces.MarketArgs{
				CurrencyPair: average.CurrencyPair,
				Interval:     int(intervel),
				Level:        average.Level,
			}).SpotMarket()

			if markets == nil {
				return 0
			}

			for i, market := range markets {
				if i == len(markets)-1 {
					continue
				}
				if i >= 1 {
					sum = sum + utils.StringToFloat64(market[2])
				}
			}

			return sum / float64(len(markets)-2)
		}
	} else {
		return 0
	}
}
