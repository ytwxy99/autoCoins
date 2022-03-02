package trade

import (
	"gorm.io/gorm"

	"github.com/gateio/gateapi-go/v6"
	"github.com/sirupsen/logrus"

	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/policy"
)

var target policy.Policy

// find macd buy point target
func FindMacdTarget(client *gateapi.APIClient, db *gorm.DB, coins []string, buyCoins chan<- string) {
	target = &policy.MacdPolicy{}
	for {
		for _, coin := range coins {
			macdCondition := target.Target(client, coin).(bool)
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

// find buy point target by doing cointegration
func DoCointegration(client *gateapi.APIClient, db *gorm.DB, buyCoins chan<- string) {
	target = &policy.Cointegration{}
	for {
		coins := target.Target(client, db).([]string)
		if len(coins) != 0 {
			for _, coin := range coins {
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
	}

}
