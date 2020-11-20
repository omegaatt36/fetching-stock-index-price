package stockindex

import "time"

// OpenOrClose defines opening or closing
type OpenOrClose string

const (
	// OpeningStr defines opening string const
	OpeningStr OpenOrClose = "opening"
	// ClosingStr defines closing string const
	ClosingStr OpenOrClose = "closing"
)

// StockPriceIndex defines stock price index history
type StockPriceIndex struct {
	Type  OpenOrClose `json:"type"`
	Date  time.Time   `json:"date"`
	Price float64     `json:"price"`
}
