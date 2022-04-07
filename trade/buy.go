package trade

import (
	"github.com/ytwxy99/autoCoins/utils/index"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"

	"github.com/ytwxy99/autoCoins/configuration"
	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/policy"
	"github.com/ytwxy99/autoCoins/utils"
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

					// fetch the 21 interval average of 4h
					averageArgs := index.Average{
						CurrencyPair: utils.IndexCoin,
						Level:        utils.Level4Hour,
						MA:           utils.MA21,
					}
					if averageArgs.Average(false) > averageArgs.Average(false) {
						body = "建议" + utils.Up + "币种: " + coin
					} else {
						body = "建议" + utils.Down + "币种: " + coin
					}
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
	//var body string
	//var sends []string
	target = &policy.Umbrella{}
	target.Target(db, sysConf)
	//i := 0
	//for i < 1 {
	//	coins := target.Target(db).([]string)
	//	if len(coins) != 0 {
	//		for _, coin := range coins {
	//			//NOTE(ytwxy99), do real trade.
	//			inOrder := database.InOrder{
	//				Contract:  coin,
	//				Direction: "up",
	//			}
	//
	//			record, err := inOrder.FetchOneInOrder(db)
	//			if err != nil {
	//				logrus.Info("Can't find in inOrder record, then will be traded : ", coin)
	//			}
	//
	//			if record == nil {
	//				// do buy
	//				err = (&inOrder).AddInOrder(db)
	//				if err != nil {
	//					logrus.Errorf("add InOrder error : %v , inOrder is %v:", err, inOrder)
	//					continue
	//				}
	//				buyCoins <- coin
	//				logrus.Info("Find it!  to get it : ", coin)
	//				sends = append(sends, coin)
	//			}
	//		}
	//
	//		for _, send := range sends {
	//			body = body + send + " "
	//		}
	//		err := utils.SendMail(sysConf, "BTC单边协整性策略", "关注币种: "+body)
	//		if err != nil {
	//			logrus.Error("Send email failed. the err is ", err)
	//		}
	//	}
	//}
}
