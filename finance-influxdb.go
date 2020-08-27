package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client"
	_ "github.com/joho/godotenv/autoload"
	"github.com/piquette/finance-go/equity"
)

// Uses https://piquette.io/projects/finance-go/ from https://github.com/piquette/finance-go - Thanks piquette

func main() {
	// Setup
	influxEndpoint := os.Getenv("INFLUX_ENDPOINT")
	if influxEndpoint == "" {
		influxEndpoint = "http://127.0.0.1:8086"
	}

	host, err := url.Parse(influxEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	clientConfig := client.Config{
		URL:      *host,
		Username: os.Getenv("INFLUX_USER"),
		Password: os.Getenv("INFLUX_PASSWORD"),
	}

	con, err := client.NewClient(clientConfig)
	if err != nil {
		log.Fatal(err)
	}
	dur, ver, err := con.Ping()
	if err != nil {
		log.Fatal(err)
	}
	if os.Getenv("DEBUG") != "" {
		fmt.Printf("Influx seems happy as a hippo! %v, %s\n", dur, ver)
	}

	// Setup collection interval
	ci, _ := strconv.ParseUint(os.Getenv("COLLECTION_INTERVAL"), 10, 32)
	if ci == 0 {
		ci = 60
	}
	collectionInterval := time.Duration(ci) * time.Second

	symbols := strings.Split(os.Getenv("SYMBOLS"), ",")
	if len(symbols) == 0 || symbols[0] == "" {
		panic("Please set the environment variable SYMBOLS to your preferred stock ticker names, comma seperated. Ie SYMBOLS=AAPL,TSLA")
	}

	data := make([]client.Point, len(symbols))

	// Collection loop
	for {
		iter := equity.List(symbols)

		// Iterate over results. Will exit upon any error.
		for iter.Next() {
			q := iter.Equity()
			if os.Getenv("DEBUG") != "" {
				fmt.Printf("%s (%s): Bid: %.2f Ask: %.2f Price: %.2f High: %.2f Low: %.2f Close: %.2f Post: %.2f Currency: %s Market State: %s\n",
					q.Symbol,
					q.ShortName,
					q.Bid,
					q.Ask,
					q.RegularMarketPrice,
					q.RegularMarketDayHigh,
					q.RegularMarketDayLow,
					q.RegularMarketPreviousClose,
					q.RegularMarketPrice+q.PreMarketChange,
					q.CurrencyID,
					q.MarketState)
			}
			data[iter.Count()] = client.Point{
				Measurement: "stocks",
				Precision:   "s",
				Fields: map[string]interface{}{
					"bid":        q.Bid,
					"ask":        q.Ask,
					"price":      q.RegularMarketPrice,
					"high":       q.RegularMarketDayHigh,
					"low":        q.RegularMarketDayLow,
					"prev_close": q.RegularMarketPreviousClose,
					"post":       q.RegularMarketPrice + q.PreMarketChange,
				},
				Tags: map[string]string{
					"symbol":      q.Symbol,
					"currency":    q.CurrencyID,
					"marketstate": string(q.MarketState),
				},
			}
		}

		// Catch an error, if there was one.
		if iter.Err() != nil {
			// Uh-oh!
			panic(iter.Err())
		}

		bps := client.BatchPoints{
			Points:   data,
			Database: "finance",
			// RetentionPolicy: "default",
		}
		_, err := con.Write(bps)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(collectionInterval)
	}
}
