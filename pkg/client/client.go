package client

import (
	"context"

	"github.com/gateio/gateapi-go/v6"

	"github.com/ytwxy99/autoCoins/pkg/configuration"
)

var Client *gateapi.APIClient

func GetClient(apiv4 *configuration.GateAPIV4) (*gateapi.APIClient, context.Context) {
	Client = gateapi.NewAPIClient(gateapi.NewConfiguration())
	// Setting host is optional. It defaults to https://api.gateio.ws/api/v4
	// client.ChangeBasePath(config.BaseUrl)
	ctx := context.WithValue(context.Background(), gateapi.ContextGateAPIV4, gateapi.GateAPIV4{
		Key:    apiv4.Key,
		Secret: apiv4.Secret,
	})

	return Client, ctx
}
