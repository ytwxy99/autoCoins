package trade

import (
	"gorm.io/gorm"

	"github.com/gateio/gateapi-go/v6"
	"github.com/sirupsen/logrus"
	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/policy"
)

// find macd buy point target
func FindMacdTarget(client *gateapi.APIClient, db *gorm.DB, coins []string, buyCoins chan<- string) {
	for {
		for _, coin := range coins {
			macdCondition := policy.MacdTarget(client, coin, buyCoins)
			if macdCondition {
				//NOTE(ytwxy99), do real trade.
				inOrder := database.InOrder{
					Contract:  coin,
					Direction: "up",
				}

				record, err := inOrder.FetchOneInOrder(db)
				if err != nil {
					logrus.Info("Can't find : ", coin)
				}

				if record == nil {
					// do buy
					err = (&inOrder).AddInOrder(db)
					if err != nil {
						logrus.Errorf("add InOrder error : %s , inOrder is %s:", err, inOrder)
						continue
					}
					buyCoins <- coin
					logrus.Info("Find it!  to get it : ", coin)
				}
			}
		}
		// waitgroup.Done()
	}

}
