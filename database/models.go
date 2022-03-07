package database

import (
	"gorm.io/gorm"
	"time"
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
	return "order"
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
	return "sold"
}

type InOrder struct {
	gorm.Model

	Contract  string
	Direction string
	Pair      string
}

func (inorder InOrder) TableName() string {
	return "inorder"
}

type HistoryDay struct {
	Contract string    `json:"contract"  gorm:"primary_key"`
	Time     time.Time `json:"time"  gorm:"primary_key"`
	Price    string
}

func (HistoryDay HistoryDay) TableName() string {
	return "history_day"
}

type Cointegration struct {
	Pair   string `json:"pair"  gorm:"primary_key"`
	Pvalue string
}

func (cointegration Cointegration) TableName() string {
	return "cointegration"
}
