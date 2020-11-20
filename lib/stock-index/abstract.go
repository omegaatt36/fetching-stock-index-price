package stockindex

import (
	"fmt"
	"time"
)

// API defines abstract interface about stock price index apis of each country
type API interface {
	Fetch(t time.Time) ([]StockPriceIndex, error)
	LoadTimeLocation() error
	GetCurrency() string
	GetTimeLocation() *time.Location
	GetOpeningTime(year, month, day int) time.Time
	GetClosingTime(year, month, day int) time.Time
}

type dataFetcherFactory interface {
	create() (API, error)
}

type attribute struct {
	timezone *time.Location
	currency string
}

// GetAPI returns stock api interface
func GetAPI(index Index) (API, error) {
	var factory dataFetcherFactory
	switch index {
	case IndexTWSE:
		factory = twseFactory{}
	case IndexDJI:
		factory = djiFactory{}
	case IndexNASDAQ:
		factory = nasdaqFactory{}
	case IndexSP500:
		factory = sp500Factory{}
	default:
		return nil, fmt.Errorf("illegal index(%d)", index)
	}

	api, err := factory.create()
	if err != nil {
		return nil, err
	}
	return api, nil
}
