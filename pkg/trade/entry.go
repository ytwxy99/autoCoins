package trade

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/ytwxy99/autocoins/database"
	"github.com/ytwxy99/autocoins/pkg/utils"
)

var entries []Entry

type Entry interface {
	PolicyEntry(ctx context.Context)
}

type Trend struct {
}

func (*Trend) PolicyEntry(ctx context.Context) {
	var buyCoins = make(chan map[string]string, 2)

	entryInit()
	go (&TrendTarget{}).DoTarget(ctx, utils.Trend, buyCoins)

	for {
		select {
		case coin := <-buyCoins:
			for cn, direction := range coin {
				logrus.Info("buy point : ", coin)
				order := database.Order{
					Contract:  cn,
					Direction: direction,
				}
				c, err := order.FetchOneOrder(ctx)
				if c == nil && err != nil {
					// buy it.
					go DoTrade(ctx, cn, direction, utils.Trend)
				}
			}
		}
	}
}

type Trend30M struct {
}

func (*Trend30M) PolicyEntry(ctx context.Context) {
	var buyCoins = make(chan map[string]string, 2)

	entryInit()
	go (&Trend30mTarget{}).DoTarget(ctx, utils.Trend30Min, buyCoins)

	for {
		select {
		case coin := <-buyCoins:
			for cn, direction := range coin {
				logrus.Info("buy point : ", coin)
				order := database.Order{
					Contract:  cn,
					Direction: direction,
				}
				c, err := order.FetchOneOrder(ctx)
				if c == nil && err != nil {
					// buy it.
					go DoTrade(ctx, cn, direction, utils.Trend30Min)
				}
			}
		}
	}
}

func entryInit() {
	// use all cpus
	runtime.GOMAXPROCS(runtime.NumCPU())

	// set pprof service
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
}
