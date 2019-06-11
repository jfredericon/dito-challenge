package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

// Representation Event Buy
type Buy struct {
	Event      string       `json:"event"`
	Timestamp  time.Time    `json:"timestamp"`
	Revenue    float64      `json:"revenue"`
	CustomData []CustomData `json:"custom_data"`
}

// Representation Event Buy Product
type BuyProduct struct {
	Event      string       `json:"event"`
	Timestamp  time.Time    `json:"timestamp"`
	CustomData []CustomData `json:"custom_data"`
}

// Representation Custom Data
type CustomData struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// Representation Timeline Item
type TimelineItem struct {
	Timestamp     time.Time `json:"timestamp"`
	Revenue       float64   `json:"revenue"`
	TransactionID string    `json:"transaction_id"`
	StoreName     string    `json:"store_name"`
	Products      []Product `json:"products"`
}

// Representation Product
type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func groupBy(maps []map[string]interface{}, key string) map[string][]map[string]interface{} {

	groups := make(map[string][]map[string]interface{})
	for _, m := range maps {
		k := m[key].(string)
		groups[k] = append(groups[k], m)
	}
	return groups

}

func sortByEventsBuy(maps []map[string]interface{}) []Buy {

	var events []Buy

	byteData, _ := json.Marshal(maps)

	json.Unmarshal(byteData, &events)

	sort.Slice(events, func(i, j int) bool { return events[j].Timestamp.Before(events[i].Timestamp) })

	return events

}

func sortByEventsBuyProduct(maps []map[string]interface{}) []BuyProduct {

	var events []BuyProduct

	byteData, _ := json.Marshal(maps)

	json.Unmarshal(byteData, &events)

	sort.Slice(events, func(i, j int) bool { return events[j].Timestamp.Before(events[i].Timestamp) })

	return events

}

func getCustomData(cd []CustomData, key string) interface{} {

	var toReturn interface{}
	for _, value := range cd {
		if value.Key == key {
			toReturn = value.Value
		}
	}

	return toReturn
}

func getProducts(transactionID interface{}, events []BuyProduct) []BuyProduct {

	var toReturn []BuyProduct

	for _, event := range events {

		if getCustomData(event.CustomData, "transaction_id") == transactionID {
			toReturn = append(toReturn, event)
		}
	}

	return toReturn

}

func createTimelineItem(eventBuy Buy, eventsBuyProduct []BuyProduct) TimelineItem {

	item := TimelineItem{
		Timestamp:     eventBuy.Timestamp,
		Revenue:       eventBuy.Revenue,
		TransactionID: getCustomData(eventBuy.CustomData, "transaction_id").(string),
		StoreName:     getCustomData(eventBuy.CustomData, "store_name").(string),
		Products:      addProducts(eventsBuyProduct),
	}

	return item

}

func addProducts(eventsBuyProduct []BuyProduct) []Product {

	var toReturn []Product

	for _, event := range eventsBuyProduct {

		product := Product{
			Name:  getCustomData(event.CustomData, "product_name").(string),
			Price: getCustomData(event.CustomData, "product_price").(float64),
		}

		toReturn = append(toReturn, product)

	}

	return toReturn
}

// Timeline : GET "/timeline" endpoint
func Timeline(netClient *http.Client) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		response, err := netClient.Get("https://storage.googleapis.com/dito-questions/events.json")

		if err != nil {
			panic(err.Error())
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)

		if err != nil {
			panic(err.Error())
		}

		dataJSON := make(map[string][]map[string]interface{})

		jsonErr := json.Unmarshal(body, &dataJSON)

		if jsonErr != nil {
			panic(jsonErr.Error())
		}

		grouped := groupBy(dataJSON["events"], "event")

		eventsBuy := grouped["comprou"]

		sortedBuy := sortByEventsBuy(eventsBuy)

		eventsBuyProduct := grouped["comprou-produto"]

		sortedBuyProduct := sortByEventsBuyProduct(eventsBuyProduct)

		var timeline []TimelineItem

		for _, event := range sortedBuy {

			transactionID := getCustomData(event.CustomData, "transaction_id")

			if transactionID != nil {
				products := getProducts(transactionID, sortedBuyProduct)

				item := createTimelineItem(event, products)

				timeline = append(timeline, item)
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"timeline": timeline,
		})
	}
}
