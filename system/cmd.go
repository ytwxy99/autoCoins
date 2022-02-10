package system

import (
	"context"
	"gorm.io/gorm"
	"time"

	"github.com/gateio/gateapi-go/v6"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	c "github.com/ytwxy99/autoCoins/client"
	"github.com/ytwxy99/autoCoins/configuration"
	"github.com/ytwxy99/autoCoins/gateway"
	"github.com/ytwxy99/autoCoins/trade"
)

// refer: https://github.com/spf13/cobra/blob/v1.2.1/user_guide.md
func InitCmd(client *gateapi.APIClient, ctx context.Context, sysConf *configuration.SystemConf, db *gorm.DB) {
	// init action
	var InitCmd = &cobra.Command{
		Use:   "init [string to echo]",
		Short: "Init trade envirment",
		Run: func(cmd *cobra.Command, args []string) {
			initErr := make(chan error)
			go func() {
				for {
					logrus.Info("update sport all coins into csv ……")
					result, err := c.GetSpotAllCoins(client, ctx)
					if err != nil {
						logrus.Error("get sport all coins error: %s\n", err)
					}
					err = InitCurrencyPairs(client, result, sysConf.CoinCsv, db)
					if err != nil {
						initErr <- err
					}

					logrus.Info("update sport all coins into csv finished!")
					// update coins list over specified interval time.
					time.Sleep(3600 * time.Second)
				}
			}()

			select {
			case err := <-initErr:
				{
					logrus.Error("update sport all coins into csv error: %s\n", err)
				}
			}
		},
	}

	var GateWayCmd = &cobra.Command{
		Use:   "gateway [string to echo]",
		Short: "Start autoCoins gateway",
		Run: func(cmd *cobra.Command, args []string) {
			router := gin.Default()
			gateway.Router(router, sysConf, db)
		},
	}

	var tradeCmd = &cobra.Command{
		Use:   "trade [string to echo]",
		Short: "Do a trade which you can choose",
		Args:  cobra.MinimumNArgs(1),
	}

	// use macd policy
	var macdCmd = &cobra.Command{
		Use:   "macd [string to echo]",
		Short: "Using macd to do a trade",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Info("market quotations is comming ！ get it !")
			for {
				t := &trade.Trade{
					Policy: "macd",
				}
				t.Entry(client, db, sysConf)
			}

		},
	}

	// use cointegration policy
	var cointegrationCmd = &cobra.Command{
		Use:   "cointegration [string to echo]",
		Short: "Using cointegration to do a trade",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Info("Find the cointegration in the sea ！ get it !")
			t := &trade.Trade{
				Policy: "cointegration",
			}
			t.Entry(client, db, sysConf)
		},
	}

	var rootCmd = &cobra.Command{Use: "autoCoin"}
	rootCmd.AddCommand(InitCmd)
	rootCmd.AddCommand(GateWayCmd)
	rootCmd.AddCommand(tradeCmd)
	tradeCmd.AddCommand(macdCmd)
	tradeCmd.AddCommand(cointegrationCmd)
	rootCmd.Execute()
}
