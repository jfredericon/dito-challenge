package main

import (
	"bytes"
	"ditoChallenge/handler"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	esURL     = "http://localhost:9200/"
	indexName = "dito_chellenge"
)

func newHTTPClient() *http.Client {
	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	netClient := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	return netClient
}

func checkElasticInstance(netClient *http.Client) {

	if _, err := netClient.Get(esURL); err != nil {
		log.Println(err)
	}
	log.Println("checkElasticInstance OK")
}

func createIndexIfNotExist(netClient *http.Client) {

	url := esURL + indexName
	log.Println(url)
	response, err := netClient.Head(url)

	if err != nil {
		log.Println(err)
	}

	if response.StatusCode == http.StatusNotFound {
		log.Println("Index Not Exist")

		indexJSON := []byte(`
		{
			"settings": {
			  "index": {
				"analysis": {
				  "filter": {},
				  "analyzer": {
					"analyzer_keyword": {
					  "tokenizer": "keyword",
					  "filter": "lowercase"
					},
					"edge_ngram_analyzer": {
					  "filter": [
						"lowercase"
					  ],
					  "tokenizer": "edge_ngram_tokenizer"
					}
				  },
				  "tokenizer": {
					"edge_ngram_tokenizer": {
					  "type": "edge_ngram",
					  "min_gram": 2,
					  "max_gram": 5,
					  "token_chars": [
						"letter"
					  ]
					}
				  }
				}
			  }
			},
			"mappings": {
				"properties": {
				  "event": {
					"type": "text",
					"analyzer": "edge_ngram_analyzer"
				  },
				  "timestamp": {
					 "type": "date"
				  }
				}
			  } 
		  }
		`)

		request, err := http.NewRequest("PUT", url, bytes.NewBuffer(indexJSON))

		if err != nil {
			log.Println(err)
		}

		request.Header.Set("Content-Type", "application/json")

		response, err := netClient.Do(request)

		if err != nil {
			log.Println(err)
		}

		log.Println(response)
		log.Println("Create Index")
	}

	log.Println("createIndexIfNotExist OK")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())

	netClient := newHTTPClient()

	checkElasticInstance(netClient)

	createIndexIfNotExist(netClient)

	v1 := r.Group("/v1")
	{
		v1.GET("/", handler.Home())
		v1.POST("/event", handler.EventCreate(netClient))
		v1.GET("/autocomplete", handler.Autocomplete(netClient))
		v1.GET("/timeline", handler.Timeline(netClient))
	}

	r.Run(":5000")
}
