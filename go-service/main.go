package main

import (
	"github.com/SwanHtetAungPhyo/stock-aggretator/src/aggregation"
	"time"
)

func main() {
	timer := time.NewTicker(2 * time.Second)
	defer timer.Stop()

	for range timer.C {
		aggregator := aggregation.NewAggregator("AAPL")
		go aggregator.PolygonRestClient()

	}
}
