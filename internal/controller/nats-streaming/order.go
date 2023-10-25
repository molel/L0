package nats_streaming

import (
	"L0/internal/entity"
	"L0/internal/repo"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
)

type OrderSubscriber struct {
	repo *repo.OrderRepo
	//Sub  stan.Subscription
}

func NewOrderSubscriber(repo *repo.OrderRepo) *OrderSubscriber {
	return &OrderSubscriber{repo: repo}
}

func (orderSub *OrderSubscriber) Subscribe(clusterId string) {
	sc, err := stan.Connect(clusterId, "order-client", stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error occurred during connecting to nats streaming: %s", err))
	}

	_, err = sc.Subscribe("orders", func(msg *stan.Msg) {
		var unmarshalledOrder entity.UnmarshalledOrder
		if err := json.Unmarshal(msg.Data, &unmarshalledOrder); err != nil {
			log.Printf("Recieved invalid data")
			return
		}

		if err := orderSub.repo.InsertOrder(unmarshalledOrder.OrderUid, unmarshalledOrder); err != nil {
			log.Printf("Failed to insert rawOrder to db: %s\n", err)
			return
		}
	})
}
