package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/nats-io/stan.go"
)

func main() {
	clusterId := "test-cluster"
	clientId := "Publisher"
	testData, err := os.Open("model.json")
	byteData, _ := ioutil.ReadAll(testData)
	if err != nil {
		fmt.Println(err)
	}
	defer testData.Close()
	sc, err := stan.Connect(
		clusterId,
		clientId,
		stan.NatsURL(stan.DefaultNatsURL),
	)
	if err != nil {
		log.Print(err)
		return
	}


	defer sc.Close()
	ackHandler := func(ackedNuid string, err error) {
		if err != nil {
			log.Printf("Warning: error publishing msg id %s: %v\n", ackedNuid, err.Error())
		} else {
			log.Printf("Received ack for msg id %s\n", ackedNuid)
		}
	}

	nuid, err := sc.PublishAsync("test", byteData, ackHandler) // returns immediately
	if err != nil {
		log.Printf("Error publishing msg %s: %v\n", nuid, err.Error())
	}


	time.Sleep(10 * time.Second)
}
