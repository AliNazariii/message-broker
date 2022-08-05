package main

import (
	"context"
	"google.golang.org/grpc"
	"math/rand"
	api "therealbroker/api/proto/src/broker/api/proto"
	"therealbroker/pkg/log"
	"time"
)

const VUs = 100
const REQUESTS = 80000

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func Publish(client api.BrokerClient, logger *log.Logger) {
	_, err := client.Publish(context.Background(), &api.PublishRequest{
		Subject: string(letters[rand.Intn(len(letters))]),
		//Subject:           "zzzzz",
		Body:              []byte("t"),
		ExpirationSeconds: 2000,
	})
	if err != nil {
		logger.Errorln("Error publishing message: ", err)
		return
	}
}

func Fetch(client api.BrokerClient, logger *log.Logger) {
	_, err := client.Fetch(context.Background(), &api.FetchRequest{
		//Subject: string(letters[rand.Intn(len(letters))]),
		Subject: "zzzzz",
		Id: 2,
	})
	if err != nil {
		logger.Errorln("Error fetching message: ", err)
		return
	}
}

func main() {
	logger := log.NewLog("debug")

	conn, err := grpc.Dial("localhost:3606", grpc.WithInsecure())
	if err != nil {
		logger.Errorln("Error connecting to broker: ", err)
		return
	}
	defer conn.Close()

	client := api.NewBrokerClient(conn)

	for i := 0; i < VUs; i++ {
		go func() {
			for j := 0; j < REQUESTS; j++ {
				Publish(client, logger)
				//Fetch(client, logger)
			}
		}()
	}
	<-time.After(time.Minute * 10)
}
