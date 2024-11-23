package main

import (
	"context"
	"fmt"
	consumer "github.com/dranikpg/gtrs"
	"github.com/redis/go-redis/v9"
	"math/rand"
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

func main() {

	opt := &redis.Options{
		Addr: "redis:6379",
	}

	rdb := redis.NewClient(opt)

	// this is used to stop the stream .. to check error cases with ack are not send  and seeds that are not processed
	ctx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFunc()

	// Create  consumers
	// - four on group g1 of stream group-stream
	groupC1 := consumer.NewGroupConsumer[WeatherReportMessage](ctx, rdb, "g1", "c1", WeatherReportStream, "0-0")
	groupC2 := consumer.NewGroupConsumer[WeatherReportMessage](ctx, rdb, "g1", "c2", WeatherReportStream, "0-0")
	groupC3 := consumer.NewGroupConsumer[WeatherReportMessage](ctx, rdb, "g1", "c3", WeatherReportStream, "0-0")
	groupC4 := consumer.NewGroupConsumer[WeatherReportMessage](ctx, rdb, "g1", "c4", WeatherReportStream, "0-0")

	//// Do some recovery when we exit.
	defer func() {

		//// Lets see what acknowledgements were not delivered
		remC1 := groupC1.Close()
		remC2 := groupC2.Close()
		remC3 := groupC3.Close()
		remC4 := groupC4.Close()

		fmt.Printf("Those acks were not sent: \ngroupC1: %v, \ngroupC2: %v, \ngroupC3: %v, \ngroupC4: %v\n",
			remC1, remC2, remC3, remC4)

	}()

	for {
		var msg consumer.Message[WeatherReportMessage]                    // our message
		var ackTarget *consumer.GroupConsumer[WeatherReportMessage] = nil // who to send the confirmation

		select {
		// Consumers just close the stream on close or cancellation without
		// sending any cancellation errors.
		// So lets not forget checking the context ourselves
		case <-ctx.Done():
			return

		// Block simultaneously on all consumers and wait for first to respond
		case msg = <-groupC1.Chan():
			ackTarget = groupC1

		case msg = <-groupC2.Chan():
			ackTarget = groupC2

		case msg = <-groupC3.Chan():
			ackTarget = groupC3

		case msg = <-groupC4.Chan():
			ackTarget = groupC4

		}

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

			// Ack blocks only if the inner ackBuffer is full.
			// Use it only inside the loop or from another goroutine with continous error processing.
			if ackTarget != nil {
				ackTarget.Ack(msg)
			}

			//// Lets sometimes send an acknowledgement to the wrong stream.
			//// Just for demonstration purposes to get a real AckError
			if ackTarget == nil && rand.Float32() < 0.1 {
				fmt.Println("Sending bad ack :)")
				groupC1.Ack(msg)
			}
		case consumer.ReadError:
			// One of the consumers will stop. So lets stop altogether.
			fmt.Fprintf(os.Stderr, "Read error! %v Exiting...\n", msg.Err)
			return
		case consumer.AckError:
			// We can identify the failed ack by stream & id
			fmt.Printf("Ack failed %v-%v :( \n", msg.Stream, msg.ID)
		case consumer.ParseError:
			// We can do something useful with errv.Data
			fmt.Println("Parse failed: raw data: ", errv.Data)
		}
	}
}
