package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

type Search struct {
	Query struct {
		Match struct {
			Event string `json:"event"`
		} `json:"match"`
	} `json:"query"`
}

type EsReturnObject struct {
	Took    int    `json:"took"`
	TimeOut bool   `json:"timed_out"`
	Hits    Hits   `json:"hits"`
	Shards  Shards `json:"_shards"`
}

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type Hits struct {
	Total struct {
		Value    int         `json:"value"`
		Relation interface{} `json:"relation"`
	} `json:"total"`
	MaxScore interface{}    `json:"max_score"`
	Hits     []InternalHits `json:"hits"`
}

type InternalHits struct {
	Index  string  `json:"_index"`
	Type   string  `json:"_type"`
	ID     string  `json:"_id"`
	Score  float32 `json:"_score"`
	Source Source  `json:"_source"`
}

type Source struct {
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
}

func search(netClient *http.Client, word string) ([]string, error) {
	url := esURL + indexName + "/_search"
	contentType := "application/json"
	toReturn := make([]string, 0)

	var search Search
	search.Query.Match.Event = word

	dataJSON, err := json.Marshal(search)

	if err != nil {
		return toReturn, err
	}

	response, err := netClient.Post(url, contentType, bytes.NewBuffer(dataJSON))

	if err != nil {
		return toReturn, err
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	var data EsReturnObject

	jsonErr := json.Unmarshal(body, &data)

	if jsonErr != nil {
		panic(jsonErr)
	}

	for _, value := range data.Hits.Hits {
		toReturn = append(toReturn, value.Source.Event)
	}

	return removeDuplicates(toReturn), nil
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}

	for v := range elements {
		encountered[elements[v]] = true
	}

	result := []string{}
	for key, _ := range encountered {
		result = append(result, key)
	}

	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })

	return result
}

// Autocomplete : GET "/" endpoint
func Autocomplete(netClient *http.Client) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		word := ctx.Query("q")

		if len(word) < 2 {
			ctx.JSON(http.StatusOK, gin.H{"data": make([]string, 0)})
			return
		}

		response, err := search(netClient, word)

		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": response,
		})

	}

}
