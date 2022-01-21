package client

import (
	"context"
	"github.com/antihax/optional"
	"github.com/gateio/gateapi-go/v6"
	"github.com/sirupsen/logrus"
	"github.com/ytwxy99/autoCoins/configuration"
)

// fetch spot client
func GetSpotClient(apiv4 *configuration.GateAPIV4) (*gateapi.APIClient, context.Context) {
	client := gateapi.NewAPIClient(gateapi.NewConfiguration())
	// Setting host is optional. It defaults to https://api.gateio.ws/api/v4
	// client.ChangeBasePath(config.BaseUrl)
	ctx := context.WithValue(context.Background(), gateapi.ContextGateAPIV4, gateapi.GateAPIV4{
		Key:    apiv4.Key,
		Secret: apiv4.Secret,
	})

	return client, ctx
}

// get all coins
func GetSpotAllCoins(client *gateapi.APIClient, ctx context.Context) ([]gateapi.CurrencyPair, error) {
	result, _, err := client.SpotApi.ListCurrencyPairs(ctx)

	if err != nil {
		if e, ok := err.(gateapi.GateAPIError); ok {
			logrus.Error("gate api error: %s\n", e.Error())
		} else {
			logrus.Error("generic error: %s\n", err.Error())
		}

		return nil, err
	}

	return result, nil
}

// get spot Market candlesticks
func GetSpotCandlesticks(client *gateapi.APIClient, currencyPair string, from int64, to int64, interval string) [][]string {
	ctx := context.Background()
	opts := &gateapi.ListCandlesticksOpts{
		From:     optional.NewInt64(from),
		To:       optional.NewInt64(to),
		Interval: optional.NewString(interval),
	}

	result, _, err := client.SpotApi.ListCandlesticks(ctx, currencyPair, opts)
	if err != nil {
		if e, ok := err.(gateapi.GateAPIError); ok {
			logrus.Error("gate api error: %s\n", e.Error())
		} else {
			logrus.Error("generic error: %s\n", err.Error())
		}
		return nil
	}

	return result
}

// Get details of a specifc order
func GetCurrencyPair(client *gateapi.APIClient, currencyPair string) ([]gateapi.Ticker, error) {
	ctx := context.Background()

	opts := &gateapi.ListTickersOpts{
		CurrencyPair: optional.NewString(currencyPair),
	}
	result, _, err := client.SpotApi.ListTickers(ctx, opts)
	if err != nil {
		if e, ok := err.(gateapi.GateAPIError); ok {
			logrus.Error("gate api error:", e.Error())
		} else {
			logrus.Error("generic error:", err.Error())
		}
		return []gateapi.Ticker{}, err
	}
	return result, err
}
