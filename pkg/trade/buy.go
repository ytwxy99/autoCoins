package trade

import (
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"

	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/pkg/configuration"
	"github.com/ytwxy99/autoCoins/pkg/policy"
	"github.com/ytwxy99/autoCoins/pkg/utils"
	"github.com/ytwxy99/autoCoins/pkg/utils/index"
)

var target policy.Policy

func FindTrendTarget(db *gorm.DB, sysConf *configuration.SystemConf, coins []string, buyCoins chan<- string) {
	target = &policy.MacdPolicy{}
	for {
		for _, coin := range coins {
			condition := target.Target(coin).(bool)
			if condition {
				//NOTE(ytwxy99), do real trade.
				inOrder := database.InOrder{
					Contract:  coin,
					Direction: utils.DirectionUp,
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
					err = utils.SendMail(sysConf, "趋势策略", "关注币种: "+coin)
					if err != nil {
						logrus.Error("Send email failed. the err is ", err)
					}
				}
			}
		}
		// waitgroup.Done()
	}
}

func DoCointegration(db *gorm.DB, sysConf *configuration.SystemConf, buyCoins chan<- string) {
	var body string
	target = &policy.Cointegration{}
	i := 0
	for i < 1 {
		coins := target.Target(db, sysConf).([]string)
		if len(coins) != 0 {
			for _, coin := range coins {
				// fetch the 21 interval average of 4h
				inOrder := database.InOrder{
					Contract: coin,
				}

				averageArgs := index.Average{
					CurrencyPair: utils.IndexCoin,
					Level:        utils.Level4Hour,
					MA:           utils.MA21,
				}
				if averageArgs.Average(false) > averageArgs.Average(false) {
					body = utils.Up + coin
					inOrder.Direction = utils.DirectionUp
				} else {
					body = utils.Down + coin
					inOrder.Direction = utils.DirectionDown
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
					err = utils.SendMail(sysConf, utils.BtcPolicy, body)
					if err != nil {
						logrus.Error("Send email failed. the err is ", err)
					}
				}
			}
		}
	}
}

func DoUmbrella(db *gorm.DB, sysConf *configuration.SystemConf, buyCoins chan<- string) {
	var body string
	var sends []string

	i := 0
	for i < 1 {
		target = &policy.Umbrella{}
		coins := target.Target(db, sysConf).([]string)
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
					sends = append(sends, coin)
				}
			}

			for _, send := range sends {
				body = body + send + " "
			}
			err := utils.SendMail(sysConf, "BNB单边协整性策略", "关注币种: "+body)
			if err != nil {
				logrus.Error("Send email failed. the err is ", err)
			}
		}
	}
}
