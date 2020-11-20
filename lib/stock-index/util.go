package stockindex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// YahooFinanceResp defines response from https://finance.yahoo.com/quote
type YahooFinanceResp struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency             string  `json:"currency"`
				Symbol               string  `json:"symbol"`
				ExchangeName         string  `json:"exchangeName"`
				InstrumentType       string  `json:"instrumentType"`
				FirstTradeDate       int     `json:"firstTradeDate"`
				RegularMarketTime    int     `json:"regularMarketTime"`
				Gmtoffset            int     `json:"gmtoffset"`
				Timezone             string  `json:"timezone"`
				ExchangeTimezoneName string  `json:"exchangeTimezoneName"`
				RegularMarketPrice   float64 `json:"regularMarketPrice"`
				ChartPreviousClose   float64 `json:"chartPreviousClose"`
				PriceHint            int     `json:"priceHint"`
				CurrentTradingPeriod struct {
					Pre struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"pre"`
					Regular struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"regular"`
					Post struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"post"`
				} `json:"currentTradingPeriod"`
				DataGranularity string   `json:"dataGranularity"`
				Range           string   `json:"range"`
				ValidRanges     []string `json:"validRanges"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Low    []float64 `json:"low"`
					Open   []float64 `json:"open"`
					Volume []int     `json:"volume"`
					Close  []float64 `json:"close"`
					High   []float64 `json:"high"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

func fetchFromYahooFinance(symbol string, t time.Time, zone *time.Location) (
	[]StockPriceIndex, error) {
	end := time.Date(t.Year(), t.Month(), t.Day(), 18, 0, 0, 0, zone)
	start := end.AddDate(0, 0, -1).Add(-time.Hour * 7)
	url := fmt.Sprintf("https://query2.finance.yahoo.com/v8/finance/chart/%s", symbol)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err,
			"want to GET %s, can't construct new request failed", req.URL.String())
	}

	q := req.URL.Query()
	q.Add("period1", fmt.Sprintf("%d", start.Unix()))
	q.Add("period2", fmt.Sprintf("%d", end.Unix()))
	q.Add("interval", "1d")
	q.Add("includeAdjustedClose", t.Format("false"))
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "GET %s failed", req.URL.String())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "GET %s, ioutil ReadAll failed, body(%s)",
			req.URL.String(), string(body))
	}
	var r YahooFinanceResp
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, errors.Wrapf(err, "GET %s, json unmarshal failed, body(%s)",
			req.URL.String(), string(body))
	}

	if r.Chart.Error != nil {
		return nil, fmt.Errorf("GET %s, but got stat:(%v)", req.URL.String(), r.Chart.Error)
	}

	historyRecords := make([]StockPriceIndex, 0)
	result := r.Chart.Result[0]
	quote := result.Indicators.Quote[0]
	for index, unixTimestamp := range result.Timestamp {

		open := quote.Open[index]
		close := quote.Close[index]
		openDate := time.Unix(unixTimestamp, 0)
		closeDate := openDate.Add(time.Hour*6 + time.Minute*30)
		historyRecords = append(historyRecords, []StockPriceIndex{
			{
				Type:  OpeningStr,
				Date:  openDate,
				Price: open,
			},
			{
				Type:  ClosingStr,
				Date:  closeDate,
				Price: close,
			},
		}...)
	}
	return historyRecords, nil
}
