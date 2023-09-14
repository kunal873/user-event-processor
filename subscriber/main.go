package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

type RequestData struct {
	UserID  string `json:"UserID"`
	Payload string `json:"Payload"`
}

var ctx = context.Background()

var redisClient = redis.NewClient(&redis.Options{
	Addr: "redis:6379",
	//    Addr: os.Getenv("redis_url") + ":" + os.Getenv("redis_port"),
})

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
	subscriber := redisClient.Subscribe(ctx, "send-user-data")

	user_data := RequestData{}

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal([]byte(msg.Payload), &user_data); err != nil {
			panic(err)
		}
		fmt.Println("Received message from " + msg.Channel + " channel.")
		fmt.Printf("%+v\n", user_data)
	}

}
