package trade

import "math"

func SellPolicy(policy string, args ...interface{}) bool {
	if policy == "macd" {
		lastPrice := args[0].(float32)
		storedPrice := args[1].(float32)

		return math.Abs(float64((lastPrice-storedPrice)/storedPrice)) >= 15
	} else if policy == "cointegration" {
		lastPrice := args[0].(float32)
		storedPrice := args[1].(float32)

		return math.Abs(float64((lastPrice-storedPrice)/storedPrice)) >= 15
	} else {
		return false
	}
}
