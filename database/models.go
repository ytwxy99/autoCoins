package database

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model

	Contract      string
	Fee_currency  string
	Price         string
	Amount        float32
	Time          int64
	Tp            float32
	Sl            float32
	Ttp           float32
	Tsl           float32
	Text          string
	Status        string
	Typee         string
	Account       string
	Side          string
	Iceberg       string
	Left          float32
	Fee           float32
	Amount_filled float32
	Direction     string
}

func (order Order) TableName() string {
	return "orders"
}

type Sold struct {
	gorm.Model

	Contract        string
	Price           string
	Volume          float32
	Time            int64
	Profit          float32
	Relative_profit string
	Test            string
	Status          string
	Typee           string
	Account         string
	Side            string
	Iceberg         string
	Direction       string
	Text            string
	Symbol          string
}

func (sold Sold) TableName() string {
	return "solds"
}

type InOrder struct {
	gorm.Model

	Contract  string `gorm:"type:varchar(32)"`
	Direction string `gorm:"type:varchar(32)"`
	Pair      string `gorm:"type:varchar(32)"`
}

func (inOrder InOrder) TableName() string {
	return "inorders"
}

type HistoryDay struct {
	Contract string    `gorm:"primary_key;index:contract_idx;type:varchar(32)"`
	Time     time.Time `gorm:"primary_key;index:time_idx"`
	Price    string    `gorm:"type:varchar(32)"`
}

func (HistoryDay HistoryDay) TableName() string {
	return "history_day"
}

type Cointegration struct {
	Pair   string `gorm:"primary_key;type:varchar(32)"`
	Pvalue string `gorm:"type:varchar(64)"`
}

func (cointegration Cointegration) TableName() string {
	return "cointegration"
}

type TradeDetail struct {
	gorm.Model

	Contract  string `gorm:"type:varchar(32)"`
	CointPair string `gorm:"type:varchar(32)"`
}

func (tradeDetail TradeDetail) TableName() string {
	return "trade_detail"
}
