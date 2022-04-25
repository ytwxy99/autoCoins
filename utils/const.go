package utils

const (
	Now               = 0
	One               = -1
	Four              = -3
	Five              = -4 // if you want to get 5 days market data, here should be typed -4
	Six               = -5
	Seven             = -6
	Ten               = -9
	Thirty            = -30
	MA5               = 5
	MA10              = 10
	MA21              = 21
	Level15Min        = "15m"
	Level4Hour        = "4h"
	Level8Hour        = "8h"
	Level1Day         = "1d"
	Close             = "c"
	Open              = "o"
	Volume            = "v"
	DBHistoryDayUniq  = "UNIQUE constraint failed: history_day.contract, history_day.time"
	IndexCoin         = "BTC_USDT"
	IndexPlatformCoin = "BNB_USDT"
	Up                = "建议做多: "
	Down              = "建议做空: "
	BtcPolicy         = "umbrella 策略"
	DirectionUp       = "up"
	DirectionDown     = "down"
	ClearOrder        = "清仓退出"
)
