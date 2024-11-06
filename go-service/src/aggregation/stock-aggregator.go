package aggregation

import (
	"encoding/json"
	"github.com/SwanHtetAungPhyo/stock-aggretator/src/config"
	modelling "github.com/SwanHtetAungPhyo/stock-aggretator/src/models"
	"github.com/polygon-io/client-go/rest/models"
	polygonws "github.com/polygon-io/client-go/websocket"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Aggregator struct {
	Symbol    string        `json:"symbol"`
	Configure config.Config `json:"configure"`
}

func NewAggregator(symbol string) *Aggregator {
	configure, _ := config.NewConfig()

	return &Aggregator{
		Symbol:    symbol,
		Configure: *configure,
	}
}

func (a *Aggregator) PolygonWSClient() {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.JSONFormatter{})

	configure, err := config.NewConfig()
	if err != nil {
		log.Error("Failed to load config:", err)
		return
	}

	cClient, err := polygonws.New(polygonws.Config{
		APIKey: configure.APIKEY,
		Feed:   polygonws.RealTime,
		Market: polygonws.Stocks,
	})
	if err != nil {
		log.Error("Failed to initialize Polygon WebSocket client:", err)
		return
	}
	defer cClient.Close()

	if err := cClient.Connect(); err != nil {
		log.Error(err)
		return
	}
	if err := cClient.Subscribe(polygonws.StocksMinAggs, a.Symbol); err != nil {
		log.Error("Subscription error:", err)
		return
	}
	if err := cClient.Connect(); err != nil {
		log.Error("WebSocket connection error:", err)
		return
	}

	// Handle interrupt signals
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	for {
		select {
		case <-sigint:
			log.Info("Received interrupt signal, closing WebSocket connection...")
			return
		case err := <-cClient.Error():
			log.Error("WebSocket client error:", err)
			return
		case out, more := <-cClient.Output():
			if !more {
				log.Info("WebSocket connection closed.")
				return
			}
			switch data := out.(type) {
			case models.Agg:
				log.WithFields(logrus.Fields{
					"open":      data.Open,
					"high":      data.High,
					"low":       data.Low,
					"close":     data.Close,
					"volume":    data.Volume,
					"timestamp": data.Timestamp,
				}).Info("Received stock aggregate data")
			default:
				log.Warn("Received unexpected data type from WebSocket")
			}
		}
	}
}

func (a *Aggregator) PolygonRestClient() {
	//	resp, err := http.Get(fmt.Sprintf("%s?%s&apiKey=%s", a.Configure.ParameterO.BaseUrl, a.getYesterdayDate(), a.Configure.ParameterO.Adjusted))
	resp, err := http.Get("https://api.polygon.io/v2/aggs/grouped/locale/us/market/stocks/2023-01-09?adjusted=true&apiKey=9_dzOMzDpl4ir26s56IyUzeODlHBm3iB")
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: received non-OK status %d", resp.StatusCode)
	}

	var stockData modelling.StockData
	if err := json.Unmarshal(body, &stockData); err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	data, err := json.Marshal(stockData)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}
	a.KafkaSent(data)
}

func (a *Aggregator) KafkaSent(data []byte) {
	writer := a.Kafka()
	err := writer.WriteMessages(nil, kafka.Message{
		Value: data,
	})

	if err != nil {
		log.Fatalf("Error sending message to Kafka: %v", err)
	}

	defer func(writer *kafka.Writer) {
		err := writer.Close()
		if err != nil {

		}
	}(writer)

	log.Println("successfully sent to the kafka")
}

func (a *Aggregator) Kafka() *kafka.Writer {
	kafkaURL := "localhost:9093"
	topic := "stock-topic"
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaURL},
		Topic:   topic,
	})
	return writer
}

func (a *Aggregator) getYesterdayDate() string {
	currentDate := time.Now()
	yesterday := currentDate.AddDate(0, 0, -1)
	return yesterday.Format("2009-09-10")
}
