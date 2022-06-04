package system

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/ytwxy99/autocoins/pkg/client"
	"github.com/ytwxy99/autocoins/pkg/configuration"
	"github.com/ytwxy99/autocoins/pkg/gateway"
	"github.com/ytwxy99/autocoins/pkg/trade"
)

// refer: https://github.com/spf13/cobra/blob/v1.2.1/user_guide.md
func InitCmd(ctx context.Context, sysConf *configuration.SystemConf, db *gorm.DB) {
	// init action
	var InitCmd = &cobra.Command{
		Use:   "init [string to echo]",
		Short: "Init trade environment",
		Run: func(cmd *cobra.Command, args []string) {
			initErr := make(chan error)
			go func() {
				for {
					logrus.Info("Initialize trading system ……")
					err := InitFutures(ctx, sysConf.UmbrellaCsv)
					if err != nil {
						logrus.Error("get all futures error: %v\n", err)
					}

					result, err := client.GetSpotAllCoins(ctx)
					if err != nil {
						logrus.Error("get all spot coins error: %v\n", err)
					}

					err = InitTrendPairs(result, sysConf.TrendCsv, db)
					if err != nil {
						initErr <- err
					}

					// use futures to statistics cointegration
					err = InitCointegrationPairs(result, sysConf.CointCsv, db)
					if err != nil {
						initErr <- err
					}

					logrus.Info("update all spot coins into csv finished!")

					err = InitCointegration(sysConf)
					if err != nil {
						initErr <- err
					}
					logrus.Info("Calculate cointegration successful!")

					// update coins list over specified interval time.
					time.Sleep(3600 * 24 * time.Second)
				}
			}()

			select {
			case err := <-initErr:
				{
					logrus.Error("Initialize trading system error: %v\n", err)
				}
			}
		},
	}

	var GateWayCmd = &cobra.Command{
		Use:   "gateway [string to echo]",
		Short: "Start autoCoins gateway",
		Run: func(cmd *cobra.Command, args []string) {
			router := gin.Default()
			gateway.Router(client.Client, router, sysConf, db)
		},
	}

	var tradeCmd = &cobra.Command{
		Use:   "trade [string to echo]",
		Short: "Do a trade which you can choose",
		Args:  cobra.MinimumNArgs(1),
	}

	// use trend policy
	var trendCmd = &cobra.Command{
		Use:   "trend [string to echo]",
		Short: "Using trend to do a trade",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Info("market quotations is comming ！ get it !")

			for {
				t := &trade.Trade{
					Policy: "trend",
				}
				t.Entry(db, sysConf)
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
			t.Entry(db, sysConf)
		},
	}

	// use cointegration policy
	var umbrellaCmd = &cobra.Command{
		Use:   "umbrella [string to echo]",
		Short: "Using umbrella to do a trade",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Info("Find the umbrella in the sea ！ get it !")
			t := &trade.Trade{
				Policy: "umbrella",
			}
			t.Entry(db, sysConf)
		},
	}

	// use trend policy
	var trend30mCmd = &cobra.Command{
		Use:   "trend30min [string to echo]",
		Short: "Using trend to do a trade",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Info("market quotations is comming ！ get it !")
			t := &trade.Trade{
				Policy: "trend30m",
			}
			t.Entry(db, sysConf)
		},
	}

	var rootCmd = &cobra.Command{Use: "autoCoin"}
	rootCmd.AddCommand(InitCmd)
	rootCmd.AddCommand(GateWayCmd)
	rootCmd.AddCommand(tradeCmd)
	tradeCmd.AddCommand(trendCmd)
	tradeCmd.AddCommand(cointegrationCmd)
	tradeCmd.AddCommand(umbrellaCmd)
	tradeCmd.AddCommand(trend30mCmd)
	rootCmd.Execute()
}
