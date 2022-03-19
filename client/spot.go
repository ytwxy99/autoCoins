package client

import (
	"context"

	"github.com/antihax/optional"
	"github.com/gateio/gateapi-go/v6"
	"github.com/sirupsen/logrus"

	"github.com/ytwxy99/autoCoins/configuration"
)

var SpotClient *gateapi.APIClient

// fetch spot client
func GetSpotClient(apiv4 *configuration.GateAPIV4) (*gateapi.APIClient, context.Context) {
	SpotClient = gateapi.NewAPIClient(gateapi.NewConfiguration())
	// Setting host is optional. It defaults to https://api.gateio.ws/api/v4
	// client.ChangeBasePath(config.BaseUrl)
	ctx := context.WithValue(context.Background(), gateapi.ContextGateAPIV4, gateapi.GateAPIV4{
		Key:    apiv4.Key,
		Secret: apiv4.Secret,
	})

	return SpotClient, ctx
}

// get all coins
func GetSpotAllCoins(ctx context.Context) ([]gateapi.CurrencyPair, error) {
	result, _, err := SpotClient.SpotApi.ListCurrencyPairs(ctx)

	if err != nil {
		if e, ok := err.(gateapi.GateAPIError); ok {
			logrus.Errorf("gate api error: %+v\n", e.Error())
		} else {
			logrus.Errorf("generic error: %+v\n", err.Error())
		}

		return nil, err
	}

	return result, nil
}

// get spot Market candlesticks
func GetSpotCandlesticks(currencyPair string, from int64, to int64, interval string) [][]string {
	ctx := context.Background()
	opts := &gateapi.ListCandlesticksOpts{
		From:     optional.NewInt64(from),
		To:       optional.NewInt64(to),
		Interval: optional.NewString(interval),
	}

	result, _, err := SpotClient.SpotApi.ListCandlesticks(ctx, currencyPair, opts)
	if err != nil {
		if e, ok := err.(gateapi.GateAPIError); ok {
			logrus.Errorf("gate api error: %+v\n", e.Error())
		} else {
			logrus.Errorf("generic error: %+v\n", err.Error())
		}
		return nil
	}

	return result
}

// Get details of a specifc order
func GetCurrencyPair(currencyPair string) ([]gateapi.Ticker, error) {
	ctx := context.Background()

	opts := &gateapi.ListTickersOpts{
		CurrencyPair: optional.NewString(currencyPair),
	}
	result, _, err := SpotClient.SpotApi.ListTickers(ctx, opts)
	if err != nil {
		if e, ok := err.(gateapi.GateAPIError); ok {
			logrus.Errorf("gate api error: %+v", e.Error())
		} else {
			logrus.Errorf("generic error: %+v", err.Error())
		}
		return []gateapi.Ticker{}, err
	}
	return result, err
}
