package system

import (
	"context"
	"github.com/ytwxy99/autocoins/pkg/utils"
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
					err := InitFutures(ctx)
					if err != nil {
						logrus.Error("get all futures error: %v\n", err)
					}

					result, err := client.GetSpotAllCoins(ctx)
					if err != nil {
						logrus.Error("get all spot coins error: %v\n", err)
					}

					err = InitTrendPairs(ctx, result)
					if err != nil {
						initErr <- err
					}

					// use futures to statistics cointegration
					err = InitCointegrationPairs(ctx, result)
					if err != nil {
						initErr <- err
					}

					logrus.Info("update all spot coins into csv finished!")

					err = InitCointegration(ctx)
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
			gateway.Router(ctx, router)
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
					Policy: utils.Trend,
				}
				t.Entry(ctx)
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
				Policy: utils.Coint,
			}
			t.Entry(ctx)
		},
	}

	// use trend policy
	var trend30mCmd = &cobra.Command{
		Use:   "trend30min [string to echo]",
		Short: "Using trend to do a trade",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Info("market quotations is comming ！ get it !")
			t := &trade.Trade{
				Policy: utils.Trend30Min,
			}
			t.Entry(ctx)
		},
	}

	var rootCmd = &cobra.Command{Use: "autoCoin"}
	rootCmd.AddCommand(InitCmd)
	rootCmd.AddCommand(GateWayCmd)
	rootCmd.AddCommand(tradeCmd)
	tradeCmd.AddCommand(trendCmd)
	tradeCmd.AddCommand(cointegrationCmd)
	tradeCmd.AddCommand(trend30mCmd)
	rootCmd.Execute()
}
