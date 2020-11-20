package main

import (
	"log"
	"time"

	"github.com/fetching-fetching-stock-index-price/lib/influxdb"
	stockindex "github.com/fetching-fetching-stock-index-price/lib/stock-index"
)

func insertIntoInfluxDB(index stockindex.Index, records []stockindex.StockPriceIndex) error {
	client, release := influxdb.Get()
	defer release()

	name := influxdb.StockPriceIndexMeasurementName(index.String())
	for _, record := range records {
		client.AsyncWrite(name,
			map[string]string{
				"board": string(record.Type),
			},
			map[string]interface{}{
				"price": record.Price,
			},
			record.Date,
		)
	}

	client.Flush()
	err := client.Flush()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	/*
		xxx.InitFromCliFlags()
	*/

	indices := []stockindex.Index{
		stockindex.IndexTWSE,
		stockindex.IndexDJI,
		stockindex.IndexNASDAQ,
		stockindex.IndexSP500,
	}
	for _, indexIndex := range indices {
		log.Println(indexIndex.String())
		API, err := stockindex.GetAPI(indexIndex)
		if err != nil {
			log.Println(err)
			continue
		}

		t := time.Now().UTC()
		historyRecords, err := API.Fetch(t)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, record := range historyRecords {
			log.Println(record)
		}
		// err = insertIntoInfluxDB(indexIndex, historyRecords)
		// if err != nil {
		// 	log.Printf("add into influxDB failed, %v\n", err)
		// 	continue
		// }
		time.Sleep(time.Second)
	}
}
