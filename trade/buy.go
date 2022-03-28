package trade

import (
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"

	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/policy"
)

var target policy.Policy

// find macd buy point target
func FindTrendTarget(db *gorm.DB, coins []string, buyCoins chan<- string) {
	target = &policy.MacdPolicy{}
	for {
		for _, coin := range coins {
			condition := target.Target(coin).(bool)
			if condition {
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
						logrus.Errorf("add InOrder error : %v , inOrder is %v:", err, inOrder)
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
func DoCointegration(db *gorm.DB, buyCoins chan<- string) {
	target = &policy.Cointegration{}
	i := 0
	for i < 1 {
		coins := target.Target(db).([]string)
		if len(coins) != 0 {
			for _, coin := range coins {
				//NOTE(ytwxy99), do real trade.
				inOrder := database.InOrder{
					Contract:  coin,
					Direction: "up",
				}

				record, err := inOrder.FetchOneInOrder(db)
				if err != nil {
					logrus.Info("Can't find in inOrder record, then will be traded : ", coin)
				}

				if record == nil {
					// do buy
					err = (&inOrder).AddInOrder(db)
					if err != nil {
						logrus.Errorf("add InOrder error : %v , inOrder is %v:", err, inOrder)
						continue
					}
					buyCoins <- coin
					logrus.Info("Find it!  to get it : ", coin)
				}
			}
		}
	}
}
