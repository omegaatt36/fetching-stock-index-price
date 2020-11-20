package influxdb

import (
	"encoding/json"
	"fmt"
	"time"

	stockindex "github.com/fetching-fetching-stock-index-price/lib/stock-index"

	inf "github.com/influxdata/influxdb1-client/v2"
	"github.com/pkg/errors"
)

// StockPriceIndexMeasurementName returns formated measurement name.
func StockPriceIndexMeasurementName(indexName string) string {
	return fmt.Sprintf("stock_price_index_%s", indexName)
}

// StockPriceIndexRequest defines the request to query stock price indicis.
type StockPriceIndexRequest struct {
	Name  string `form:"name" binding:"required"`
	Begin string `form:"begin" binding:"required,datetime=2006-01-02"`
	End   string `form:"end" binding:"required,datetime=2006-01-02"`
}

// QueryStockPriceIndices return stock price index recors in influxdb
func (c *Client) QueryStockPriceIndices(req StockPriceIndexRequest) (
	[]stockindex.StockPriceIndex, error) {
	index, err := stockindex.ParseIndex(req.Name)
	if err != nil {
		return nil, errors.Wrap(err, req.Name)
	}
	api, err := stockindex.GetAPI(index)
	if err != nil {
		return nil, err
	}

	measurementName := StockPriceIndexMeasurementName(req.Name)
	var bY, bM, bD, eY, eM, eD int
	fmt.Sscanf(req.Begin, "%d-%d-%02d", &bY, &bM, &bD)
	fmt.Sscanf(req.End, "%d-%d-%02d", &eY, &eM, &eD)
	begin := api.GetOpeningTime(bY, bM, bD).Add(-time.Second)
	end := api.GetClosingTime(eY, eM, eD).Add(time.Second)
	sql := fmt.Sprintf("SELECT * FROM %s WHERE time >= '%v' AND time <= '%v'",
		measurementName,
		begin.Format(time.RFC3339),
		end.Format(time.RFC3339),
	)

	q := inf.NewQuery(sql, c.dbName, "")
	response, err := c.Client.Query(q)
	if err != nil {
		return nil, errors.Wrapf(err, sql)
	}
	if err = checkResponse(response); err != nil {
		return nil, errors.Wrapf(err, sql)
	}
	location := api.GetTimeLocation()
	records := make([]stockindex.StockPriceIndex, 0, len(response.Results[0].Series[0].Values))
	for _, row := range response.Results[0].Series[0].Values {
		typeS, ok1 := row[1].(string)
		dateS, ok2 := row[0].(string)
		priceN, ok3 := row[2].(json.Number)
		if !ok1 || !ok2 || !ok3 {
			continue
		}

		date, _ := time.Parse(time.RFC3339, dateS)
		price, _ := priceN.Float64()
		records = append(records, stockindex.StockPriceIndex{
			Type:  stockindex.OpenOrClose(typeS),
			Date:  date.In(location),
			Price: price,
		})
	}
	return records, nil
}
