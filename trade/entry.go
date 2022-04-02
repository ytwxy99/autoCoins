package trade

import (
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/ytwxy99/autoCoins/configuration"
	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/utils"
)

type Trade struct {
	Policy string
}

// trade entry point
func (t *Trade) Entry(db *gorm.DB, sysConf *configuration.SystemConf) {
	var buyCoins = make(chan string, 2)
	// use all cpus
	runtime.GOMAXPROCS(runtime.NumCPU())

	// set pprof service
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	if t.Policy == "trend" {
		coins, err := utils.ReadLines(sysConf.CointCsv)
		if err != nil {
			logrus.Error("Read local file error: %v", err)
			return
		}

		for i := 0; i < (len(coins)/10 + 1); i++ {
			if i == len(coins)/10 {
				go FindTrendTarget(db, sysConf, coins[i*10:i*10+len(coins)%10], buyCoins)
			} else {
				go FindTrendTarget(db, sysConf, coins[i*10:i*10+9], buyCoins)
			}
		}

		for {
			select {
			case coin := <-buyCoins:
				logrus.Info("buy point : ", coin)
				order := database.Order{
					Contract:  coin,
					Direction: "up",
				}
				c, err := order.FetchOneOrder(db)
				if c == nil && err != nil {
					// buy it.
					go DoTrade(db, sysConf, coin, "up", "trend", SellPolicy)
				}
			}
		}

	} else if t.Policy == "cointegration" {
		var buyCoins = make(chan string, 2)
		go DoCointegration(db, sysConf, buyCoins)

		for {
			select {
			case coin := <-buyCoins:
				logrus.Info("buy point : ", coin)
				order := database.Order{
					Contract:  coin,
					Direction: "up",
				}
				c, err := order.FetchOneOrder(db)
				if c == nil && err != nil {
					// buy it.
					go DoTrade(db, sysConf, coin, "up", "cointegration", SellPolicy)
				}
			}
		}
	} else if t.Policy == "umbrella" {
		var buyCoins = make(chan string, 2)
		DoUmbrella(db, sysConf, buyCoins)

		//for {
		//	select {
		//	case coin := <-buyCoins:
		//		logrus.Info("buy point : ", coin)
		//		order := database.Order{
		//			Contract:  coin,
		//			Direction: "up",
		//		}
		//		c, err := order.FetchOneOrder(db)
		//		if c == nil && err != nil {
		//			// buy it.
		//			go DoTrade(db, sysConf, coin, "up", "cointegration", SellPolicy)
		//		}
		//	}
		//}
	}
}
