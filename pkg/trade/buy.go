package trade

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/ytwxy99/autocoins/database"
	"github.com/ytwxy99/autocoins/pkg/policy"
	"github.com/ytwxy99/autocoins/pkg/utils"
)

type Target interface {
	DoTarget(ctx context.Context, policy string, buyCoins chan<- map[string]string)
}

type TrendTarget struct {
}

func (*TrendTarget) DoTarget(ctx context.Context, policyName string, buyCoins chan<- map[string]string) {
	coins, err := utils.ReadLines(utils.GetSystemConfContext(ctx).TrendCsv)
	if err != nil {
		logrus.Error("Read trend csv failed. the err is ", err)
	}

	for {
		for _, coin := range coins {
			results := (&policy.TrendPolicy{}).Target(utils.SetContextValue(ctx, "coin", coin))
			for contract, direction := range results {
				// do real trade.
				inOrder := database.InOrder{
					Contract:  contract,
					Direction: direction,
				}

				record, err := inOrder.FetchOneInOrder(ctx)
				if err != nil {
					logrus.Info("Can't find : ", coin)
				}

				if record == nil {
					// do buy
					err = (&inOrder).AddInOrder(ctx)
					if err != nil {
						logrus.Errorf("add InOrder error : %v , inOrder is %v:", err, inOrder)
						continue
					}
					buyCoins <- results
					logrus.Info("Find it!  to get it : ", coin)
					err = utils.SendMail(utils.GetSystemConfContext(ctx), "趋势策略", "关注币种: "+coin)
					if err != nil {
						logrus.Error("Send email failed. the err is ", err)
					}
				}
			}
		}
	}
}

//func DoCointegration(ctx context.Context, buyCoins chan<- string) {
//	var body string
//	target = &policy.Cointegration{}
//	i := 0
//	for i < 1 {
//		coins := target.Target(ctx).([]string)
//		if len(coins) != 0 {
//			for _, coin := range coins {
//				// fetch the 21 interval average of 4h
//				inOrder := database.InOrder{
//					Contract: coin,
//				}
//
//				averageArgs := index.Average{
//					CurrencyPair: utils.IndexCoin,
//					Level:        utils.Level4Hour,
//					MA:           utils.MA21,
//				}
//				if averageArgs.Average(false) > averageArgs.Average(false) {
//					body = utils.Up + coin
//					inOrder.Direction = utils.DirectionUp
//				} else {
//					body = utils.Down + coin
//					inOrder.Direction = utils.DirectionDown
//				}
//
//				record, err := inOrder.FetchOneInOrder(ctx)
//				if err != nil {
//					logrus.Info("Can't find in inOrder record, then will be traded : ", coin)
//				}
//
//				if record == nil {
//					// do buy
//					err = (&inOrder).AddInOrder(ctx)
//					if err != nil {
//						logrus.Errorf("add InOrder error : %v , inOrder is %v:", err, inOrder)
//						continue
//					}
//					buyCoins <- coin
//					logrus.Info("Find it!  to get it : ", coin)
//					err = utils.SendMail(utils.GetSystemConfContext(ctx), utils.BtcPolicy, body)
//					if err != nil {
//						logrus.Error("Send email failed. the err is ", err)
//					}
//				}
//			}
//		}
//	}
//}

type Trend30mTarget struct {
}

func (*Trend30mTarget) DoTarget(ctx context.Context, policyName string, buyCoins chan<- map[string]string) {

	for {
		coins, err := utils.ReadLines(utils.GetSystemConfContext(ctx).WeightCsv)
		if err != nil {
			logrus.Error("read weight csv failed, err is ", err)
		}

		for _, coin := range coins {
			target := (&policy.Trend30M{}).Target(utils.SetContextValue(ctx, "coin", coin)).(map[string]string)
			if len(target) == 0 {
				continue
			}

			// do real trade.
			inOrder := database.InOrder{
				Contract:  coin,
				Direction: target[coin],
			}

			record, err := inOrder.FetchOneInOrder(ctx)
			if err != nil {
				logrus.Info("Can't find : ", coin)
			}

			if record == nil {
				// do buy
				err = (&inOrder).AddInOrder(ctx)
				if err != nil {
					logrus.Errorf("add InOrder error : %v , inOrder is %v:", err, inOrder)
					continue
				}
				logrus.Info("Find it!  to get it : ", target)
				err = utils.SendMail(utils.GetSystemConfContext(ctx), "交易推荐", "关注币种: "+coin+" 方向："+target[coin])
				if err != nil {
					logrus.Error("Send email failed. the err is ", err)
				}
				buyCoins <- target
			}
		}
	}
}
