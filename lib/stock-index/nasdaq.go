package stockindex

import (
	"time"

	"github.com/pkg/errors"
)

type nasdaqFactory struct{}

func (nasdaqFactory) create() (API, error) {
	api := nasdaq{}
	api.attribute.currency = "USD"
	if err := api.LoadTimeLocation(); err != nil {
		return nil, err
	}
	return &api, nil
}

type nasdaq struct {
	attribute
}

func (api *nasdaq) Fetch(t time.Time) ([]StockPriceIndex, error) {
	return fetchFromYahooFinance("%5EIXIC", t, api.timezone)
}

func (api *nasdaq) LoadTimeLocation() error {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		return errors.Wrapf(err, "can't load timezone.")
	}
	api.attribute.timezone = location
	return nil
}

func (api *nasdaq) GetTimeLocation() *time.Location {
	return api.attribute.timezone
}

func (api *nasdaq) isAmericaSummerTime(year, month, day int) bool {
	parseDayNumber := func(d time.Time) int {
		return d.Year()*10000 + int(d.Month())*100 + d.Day()
	}
	summerStart := time.Date(year, 3, 1, 0, 0, 0, 0, api.attribute.timezone)
	for count := 0; count < 2; {
		if summerStart.Weekday() == time.Sunday {
			count++
		}
		summerStart = summerStart.AddDate(0, 0, 1)
	}

	summerEnd := time.Date(year, 11, 1, 0, 0, 0, 0, api.attribute.timezone)
	for count := 0; count < 1; {
		if summerEnd.Weekday() == time.Sunday {
			count++
		}
		summerEnd = summerEnd.AddDate(0, 0, 1)
	}

	if now := year*1000 + month*100 + day; now >= parseDayNumber(summerStart) &&
		now < parseDayNumber(summerEnd) {
		return true
	}
	return false
}

func (api *nasdaq) GetCurrency() string {
	return api.attribute.currency
}

func (api *nasdaq) GetOpeningTime(year, month, day int) time.Time {
	if api.isAmericaSummerTime(year, month, day) {
		return time.Date(year, time.Month(month), day, 8, 30, 0, 0, api.attribute.timezone)
	}
	return time.Date(year, time.Month(month), day, 9, 30, 0, 0, api.attribute.timezone)
}

func (api *nasdaq) GetClosingTime(year, month, day int) time.Time {
	if api.isAmericaSummerTime(year, month, day) {
		return time.Date(year, time.Month(month), day, 15, 00, 0, 0, api.attribute.timezone)
	}
	return time.Date(year, time.Month(month), day, 16, 00, 0, 0, api.attribute.timezone)
}
