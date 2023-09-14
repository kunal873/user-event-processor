package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

type RequestData struct {
	UserID  string `json:"UserID"`
	Payload string `json:"Payload"`
}

var redisClient = redis.NewClient(&redis.Options{Addr: "redis:6379"})

var ctx = context.Background()

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
	port := os.Getenv("PORT")

	r := gin.Default()

	go processEventsFromStream()

	r.POST("/api/publish", func(c *gin.Context) {
		var requestData RequestData

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID := requestData.UserID
		payload := requestData.Payload

		// Append the event to a Redis Stream
		xAdd := redis.XAddArgs{
			Stream: "user-events",
			Values: map[string]interface{}{
				"UserID":  userID,
				"Payload": payload,
			},
		}

		if err := redisClient.XAdd(ctx, &xAdd).Err(); err != nil {
			fmt.Println("Failed to publish event", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish event"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Data published successfully", "UserID": userID, "Payload": payload})
	})

	r.Run(":" + port)
}

func processEventsFromStream() {
	for {
		// Read events from the Redis Stream
		streams, err := redisClient.XRead(ctx, &redis.XReadArgs{
			Streams: []string{"user-events", "0"},
			Block:   0,
			Count:   1,
		}).Result()
		if err != nil {
			continue
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				fmt.Println("Received event with ID", message.ID)
				userID := message.Values["UserID"].(string)
				payload := message.Values["Payload"].(string)

				// send event with retry logic
				maxRetries := 3
				for retries := 0; retries < maxRetries; retries++ {
					if processEvent(userID, payload) {
						fmt.Println("user-events processed successfully")
						break
					} else {
						fmt.Println("Event processing failed, retrying...")
					}
				}
				// remove the event from the stream
				redisClient.XDel(ctx, stream.Stream, message.ID)
			}
		}
	}
}

func processEvent(userID, payload string) bool {
	payload_data := map[string]interface{}{
		"UserID":  userID,
		"Payload": payload,
	}

	payload_data_1, err := json.Marshal(payload_data)
	if err != nil {
		return false
	}
	if err := redisClient.Publish(ctx, "send-user-data", payload_data_1).Err(); err != nil {
		return false
	}
	fmt.Println("message sent to subscriber")
	return true
}
