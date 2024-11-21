package main

import (
	"context"
	"fmt"
	stream "github.com/dranikpg/gtrs"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

var WeatherReportStream = "weather-station-stream"

type WeatherReportMessage struct {
	Station     string
	Temperature int
	Humidity    int
	EventTime   string
}

var StartStreamID = "0-0"

func main() {

	// HINT: initialize redis client to connect to server
	opt := &redis.Options{
		Addr: "redis:6379",
	}
	rdb := redis.NewClient(opt)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 200*time.Second)

	defer cancelFunc()

	////TODO: Get the right
	//fmt.Println("INFO: Reading from last success")
	//lastId, _ := rdb.Get(ctx, "CONSUMER#KEY#LASTID").Result()
	//fmt.Println("Last ID: ", lastId)
	//
	//TODO: Create a client for "main-stream" redis stream
	mainStream := stream.NewConsumer[WeatherReportMessage](ctx, rdb,
		stream.StreamIDs{
			WeatherReportStream: StartStreamID},
	)

	defer func() {

		// Lets see where we stopped reading "main-stream"
		seenIds := mainStream.Close()
		fmt.Println("Main stream reader stopped on", seenIds)

		//fmt.Println("Saving last Read")
		//_, err := rdb.Set(context.Background(), "CONSUMER#KEY#LASTID", seenIds["main-stream"], 0).Result()
		//if err != nil {
		//	fmt.Println("Error:", err)
		//}
	}()

	for {
		var msg stream.Message[WeatherReportMessage] // our message

		select {
		// Consumers just close the stream on close or cancellation without
		// sending any cancellation errors.
		// So lets not forget checking the context ourselves
		case <-ctx.Done():
			return

		// Block for consumer and wait for first to respond
		case msg = <-mainStream.Chan():

			switch errv := msg.Err.(type) {

			// This interface-nil comparison in safe. Consumers never return typed nil errors.
			case nil:
				fmt.Printf("Got event %v: Station: %v, Temperature: %v, Humidity: %v, Time: %v, from stream: <%v>\n",
					msg.ID,
					msg.Data.Station,
					msg.Data.Temperature,
					msg.Data.Humidity,
					msg.Data.EventTime,
					msg.Stream)

			case stream.ReadError:
				fmt.Fprintf(os.Stderr, "Read error! %v Exiting...\n", msg.Err)
				return

			case stream.AckError:
				// We can identify the failed ack by stream & id
				fmt.Printf("Ack failed %v-%v :( \n", msg.Stream, msg.ID)

			case stream.ParseError:
				// We can do something useful with errv.Data
				fmt.Println("Parse failed: raw data: ", errv.Data)
			}
		}
	}
}
