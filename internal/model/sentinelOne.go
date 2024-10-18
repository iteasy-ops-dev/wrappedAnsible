package model

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"iteasy.wrappedAnsible/pkg/utils"
)

// Requester defines methods that must be implemented by each data model
type Requester interface {
	GetCollection() string
	GetAPIURL() string
	GetAPIKey() string
	GetFilterKey() string
}

// BaseRequester contains shared logic
type BaseRequester struct {
	Requester
}

func (br *BaseRequester) _collection() *mongo.Collection {
	return db.Collection(br.GetCollection())
}

// Update method to fetch data from API and upsert into MongoDB
func (br *BaseRequester) Update() error {
	col := br._collection()
	nextCursor := ""
	var wg sync.WaitGroup

	headers := map[string]string{
		"Authorization": fmt.Sprintf("ApiToken %s", br.GetAPIKey()),
		"Content-Type":  "application/json",
	}

	for {
		apiURL := br.GetAPIURL()
		if nextCursor != "" {
			apiURL += "&cursor=" + nextCursor
		}

		req, _ := http.NewRequest("GET", apiURL, nil)
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("❌ Error fetching data: %v", err)
			return err
		}
		defer resp.Body.Close()

		var data struct {
			Data       []map[string]interface{} `json:"data"`
			Pagination struct {
				NextCursor string `json:"nextCursor"`
			} `json:"pagination"`
		}

		if err := utils.ParseJSONBody(resp.Body, &data); err != nil {
			log.Printf("❌ Error parsing response body: %v", err)
			return err
		}

		nextCursor = data.Pagination.NextCursor

		for _, item := range data.Data {
			wg.Add(1)
			go func(item map[string]interface{}) {
				defer wg.Done()
				filterKey := br.GetFilterKey()
				filter := bson.M{filterKey: item[filterKey]}
				update := bson.M{"$set": item}
				_, err := col.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
				if err != nil {
					log.Printf("❌ Failed to upsert %s %v: %v", filterKey, item[filterKey], err)
				}
			}(item)
		}

		if nextCursor == "" {
			break
		}
	}

	wg.Wait()
	return nil
}
