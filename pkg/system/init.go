package system

import (
	"context"
	"os/exec"
	"time"

	"github.com/gateio/gateapi-go/v6"
	"github.com/sirupsen/logrus"
	"github.com/ytwxy99/autocoins/database"
	"github.com/ytwxy99/autocoins/pkg/client"
	"github.com/ytwxy99/autocoins/pkg/interfaces"
	"github.com/ytwxy99/autocoins/pkg/utils"
)

// InitTrendPairs fetch trend coins
func InitTrendPairs(ctx context.Context, pairs []gateapi.CurrencyPair) error {
	coins := []string{}
	for _, pair := range pairs {
		// just record coin which is tradable
		if pair.TradeStatus == "tradable" && pair.Quote == "USDT" {
			values := (&interfaces.MarketArgs{
				CurrencyPair: pair.Id,
				Interval:     -150,
				Level:        utils.Level1Day,
			}).SpotMarket()

			if len(values) < 150 {
				continue
			}

			if len(values) != 0 && utils.StringToFloat32(values[0][1]) >= 50000.0 {
				coins = append(coins, pair.Id)
			}
		}
	}

	return utils.WriteLines(coins, utils.GetSystemConfContext(ctx).TrendCsv)
}

func InitCointegrationPairs(ctx context.Context, pairs []gateapi.CurrencyPair) error {
	var historyDay *database.HistoryDay
	coins := []string{}
	for _, pair := range pairs {
		// just record coin which is tradable
		if pair.TradeStatus == "tradable" && pair.Quote == "USDT" {
			values := (&interfaces.MarketArgs{
				CurrencyPair: pair.Id,
				Interval:     -900,
				Level:        utils.Level1Day,
			}).SpotMarket()

			if len(values) < 900 {
				continue
			}

			for _, value := range values {
				timeTrans, err := time.ParseInLocation("2006-01-02 08:00:00", value[0], time.Local)
				if err != nil {
					logrus.Error("get time type from string error: %v\n", err)
				}

				historyDay = &database.HistoryDay{
					Contract: pair.Id,
					Time:     timeTrans,
					Price:    value[2],
				}

				err = historyDay.AddHistoryDay(utils.GetDBContext(ctx))
				if err != nil {
					if err.Error() != utils.DBHistoryDayUniq {
						logrus.Error("add HistoryDay record error: %v\n", err)
					}
				}
			}

			if len(values) != 0 {
				coins = append(coins, pair.Id)
			}
		}
	}

	return utils.WriteLines(coins, utils.GetSystemConfContext(ctx).CointCsv)
}

func InitCointegration(ctx context.Context) error {
	sysConf := utils.GetSystemConfContext(ctx)
	args := []string{
		sysConf.CointegrationSrcipt, // index 0
		sysConf.CointCsv,            // index 1
		sysConf.DBType,              // index 2
		sysConf.DBPath,              // index 3
		sysConf.Mysql.User,
		sysConf.Mysql.Password,
		sysConf.Mysql.Port,
		sysConf.Mysql.Host,
		sysConf.Mysql.Database,
	}
	cmd := exec.Command("python3", args...)

	_, err := cmd.Output()
	if err != nil {
		logrus.Error("run cointegration python srcipt failed. ", "error: ", err)
		return err
	}

	return nil
}

// InitFutures init future coins into Umbrella.csv
func InitFutures(ctx context.Context) error {
	coins := []string{}
	futures, err := (&interfaces.Future{
		Settle: client.Settle,
	}).GetAllFutures(ctx)
	if err != nil {
		logrus.Error("get all futures failed.", err)
	}

	for _, future := range futures {
		coins = append(coins, future.Name)
	}

	return utils.WriteLines(coins, utils.GetSystemConfContext(ctx).UmbrellaCsv)
}

//Init system base data
func Init() {
	authConf, _ := utils.ReadGateAPIV4("./etc/auth.yml")
	sysConf, _ := utils.ReadSystemConfig("./etc/autoCoin.yml")
	_, ctx := client.GetClient(authConf)
	utils.InitLog(sysConf.LogPath)
	db := database.GetDB(sysConf)
	database.InitDB(db)

	ctxMetadata := utils.SystemContext{
		SystemConf: sysConf,
		Database:   db,
	}
	InitCmd(utils.SetContextValue(ctx, "ctxMetadata", ctxMetadata), sysConf, db)
}
