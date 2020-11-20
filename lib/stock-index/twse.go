package stockindex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type twseFactory struct{}

func (twseFactory) create() (API, error) {
	api := twse{}
	api.attribute.currency = "NTD"
	if err := api.LoadTimeLocation(); err != nil {
		return nil, err
	}
	return &api, nil
}

type twse struct {
	attribute
}

func (api *twse) Fetch(t time.Time) ([]StockPriceIndex, error) {
	type Resp struct {
		Stat string     `json:"stat"`
		Data [][]string `json:"data"`
	}

	req, err := http.NewRequest("GET", "https://www.twse.com.tw/indicesReport/MI_5MINS_HIST", nil)
	if err != nil {
		return nil, errors.Wrapf(err,
			"want to GET %s, can't construct new request failed", req.URL.String())
	}

	q := req.URL.Query()
	q.Add("date", t.Format("20060102"))
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
	var r Resp
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, errors.Wrapf(err, "GET %s, json unmarshal failed, body(%s)",
			req.URL.String(), string(body))
	}

	if r.Stat != "OK" {
		return nil, fmt.Errorf("GET %s, but got stat:(%s)", req.URL.String(), r.Stat)
	}

	parseFloat := func(s string) (float64, error) {
		s = strings.ReplaceAll(strings.Trim(s, " "), ",", "")
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		return f, nil
	}

	historyRecords := make([]StockPriceIndex, 0)
	for _, d := range r.Data {
		// year need plus 1911 to AD year
		var year, month, day int
		count, err := fmt.Sscanf(d[0], "%d/%d/%02d", &year, &month, &day)
		if err != nil || count != 3 {
			return nil, errors.Wrapf(err, "date(%s) format error, want 3 but got %d",
				d[0], count)
		}
		priceList := make([]float64, 2)
		for index, priceStr := range []string{d[1], d[4]} {
			price, err := parseFloat(priceStr)
			if err != nil {
				return nil, errors.Wrapf(err, "price(%s) parse failed", priceStr)
			}
			priceList[index] = price
		}
		historyRecords = append(historyRecords, []StockPriceIndex{
			{
				Type:  OpeningStr,
				Date:  api.GetOpeningTime(year+1911, month, day),
				Price: priceList[0],
			},
			{
				Type:  ClosingStr,
				Date:  api.GetClosingTime(year+1911, month, day),
				Price: priceList[1],
			},
		}...)
	}
	return historyRecords, nil
}

func (api *twse) LoadTimeLocation() error {
	location, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return errors.Wrapf(err, "can't load timezone.")
	}
	api.attribute.timezone = location
	return nil
}

func (api *twse) GetCurrency() string {
	return api.attribute.currency
}
func (api *twse) GetTimeLocation() *time.Location {
	return api.attribute.timezone
}

func (api *twse) GetOpeningTime(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 9, 30, 0, 0, api.attribute.timezone)
}

func (api *twse) GetClosingTime(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 13, 30, 0, 0, api.attribute.timezone)
}
