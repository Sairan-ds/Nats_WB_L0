package streaming

import (
	"encoding/json"

	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"github.com/sairan-ds/go-nats-steaming-project/internal/database"
)

type Subscriber struct {
	Sc  *stan.Conn
	Sub stan.Subscription
	Db  *database.Database
}
// Connection to Nats Streaming server
func Connect() *Subscriber {
	clusterId := os.Getenv("NATS_CLUSTER_ID")
	clientId := os.Getenv("NATS_CLIENT_ID")
	sub := Subscriber{}
	sc, err := stan.Connect(
		clusterId,
		clientId,
		stan.NatsURL(stan.DefaultNatsURL),
		stan.Pings(5, 5),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("connection lost, reason: %v", reason)
		}))

	if err != nil {
		log.Fatal(err)
	}
	sub.Sc = &sc
	return &sub
	
}
// Subscribe to nats streaming server
func Subscribe(db *database.Database) *Subscriber {
	s := Connect()
	s.Db = db
	log.Println("Setting up new Subscription")
	mcb := func(msg *stan.Msg) {
		if err := msg.Ack(); err != nil {
			log.Println(err)
		}
		var order database.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Printf("Invalid data %v", err)
			return
		}

		s.Db.AddOrder(&order)
		log.Println("New order, ID:  ", order.OrderUid)
	}

	sub, err := (*s.Sc).QueueSubscribe(
		"test",
		"test",
		mcb,
		stan.SetManualAckMode())
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	s.Sub = sub
return s
}

func (s Subscriber) Unsubscribe() {
	s.Sub.Unsubscribe()
}

