package client

import (
	"context"

	"github.com/antihax/optional"
	"github.com/gateio/gateapi-go/v6"
	"github.com/sirupsen/logrus"
)

const Settle = "usdt"

func GetFutureCandlesticks(contract string, from int64, to int64, interval string) []gateapi.FuturesCandlestick {
	ctx := context.Background()
	opts := &gateapi.ListFuturesCandlesticksOpts{
		From:     optional.NewInt64(from),
		To:       optional.NewInt64(to),
		Interval: optional.NewString(interval),
	}

	resps, _, err := Client.FuturesApi.ListFuturesCandlesticks(ctx, Settle, contract, opts)
	if err != nil {
		if e, ok := err.(gateapi.GateAPIError); ok {
			logrus.Errorf("gate api error: %+v\n", e.Error())
		} else {
			logrus.Errorf("generic error: %+v\n", err.Error())
		}
		return nil
	}

	return resps
}