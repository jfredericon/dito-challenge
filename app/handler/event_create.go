package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	esURL     = "http://localhost:9200/"
	indexName = "dito_chellenge"
)

type Event struct {
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
}

func add(netClient *http.Client, newEvent Event) error {

	url := esURL + indexName + "/_doc/?refresh"
	contentType := "application/json"
	dataJSON, err := json.Marshal(newEvent)

	if err != nil {
		return err
	}

	if _, err := netClient.Post(url, contentType, bytes.NewBuffer(dataJSON)); err != nil {
		return err
	}

	return nil
}

// EventCreate :  POST "/event" endpoint
func EventCreate(netClient *http.Client) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		var body Event

		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newEvent := Event{body.Event, time.Now()}

		if err := add(netClient, newEvent); err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Event created",
			"data":    newEvent,
		})
	}
}
