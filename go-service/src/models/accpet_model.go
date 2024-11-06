package models

type Stock struct {
	T         string  `json:"T"`           // Ticker symbol
	V         float64 `json:"v"`           // Volume
	VW        float64 `json:"vw"`          // Weighted average price
	O         float64 `json:"o"`           // Opening price
	C         float64 `json:"c"`           // Closing price
	H         float64 `json:"h"`           // Highest price
	L         float64 `json:"l"`           // Lowest price
	Timestamp int64   `json:"t"`           // Time in Unix timestamp (milliseconds)
	N         int     `json:"n,omitempty"` // Number of transactions (optional)
}

// StockData represents the response structure containing multiple stocks and metadata
type StockData struct {
	QueryCount   int     `json:"queryCount"`   // Total number of queries
	ResultsCount int     `json:"resultsCount"` // Total number of results
	Adjusted     bool    `json:"adjusted"`     // Whether data is adjusted
	Results      []Stock `json:"results"`      // List of stock data
	Status       string  `json:"status"`       // Status of the request
	RequestID    string  `json:"request_id"`   // Unique request ID
	Count        int     `json:"count"`        // Count of results
}
