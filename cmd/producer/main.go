package main

import (
	"context"
	"fmt"
	stream "github.com/dranikpg/gtrs"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"time"
)

type WeatherReportMessage struct {
	Station     string
	Temperature int
	Humidity    int
	EventTime   string
}

var CityWeatherStations = []string{
	"Athens", "Helsinki", "Budapest", "London", "Paris", "Berlin", "Sofia", "Warsaw", "Rome", "Madrid"}

var WeatherReportStream = "weather-station-stream"

func main() {

	produce(10, 10, 2*time.Second, WeatherReportStream)

}

func produce(messages int, batches int, delay time.Duration, streamName string) {
	fmt.Printf("Welcome to weather-station producer. Number of messages to be send: %d\n", messages)

	// Initialize a Redis client
	opt := &redis.Options{
		Addr: "redis:6379",
	}
	ctx := context.Background()

	client := redis.NewClient(opt)

	// Crete a new Stream
	producer := stream.NewStream[WeatherReportMessage](client, streamName, nil)

	for batch := 0; batch < batches; batch++ {
		fmt.Printf("Sending batch number: %d of %d messages to stream: %v\n",
			batch, messages, streamName)
		for i := 0; i < messages; i++ {

			//TODO: START Sending messages
			_, err := producer.Add(ctx, WeatherReportMessage{
				Station:     CityWeatherStations[rand.Intn(9)],
				Temperature: rand.Intn(45),
				Humidity:    rand.Intn(100),
				EventTime:   time.Now().Format(time.ANSIC),
			})

			if err != nil {
				fmt.Printf("Unable to send message to stream: %v\n", streamName)
			}
		}
		time.Sleep(delay)
	}
}
