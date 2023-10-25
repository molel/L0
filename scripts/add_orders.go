package main

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"time"
)

type Order struct {
	OrderUid          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmId              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
	Delivery          `json:"delivery"`
	Payment           `json:"payment"`
	Items             []Item `json:"items"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func main() {
	sc, err := stan.Connect("test-cluster", "test-test", stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		log.Println(err)
	}

	validOrder := Order{
		OrderUid:          "123",
		TrackNumber:       "track numnber smtrh",
		Entry:             "entry sgfmsadg",
		Locale:            "en",
		InternalSignature: "sign",
		CustomerId:        "cid",
		DeliveryService:   "selserv",
		Shardkey:          "shke",
		SmId:              123231,
		DateCreated:       time.Now(),
		OofShard:          "oooffff",
		Delivery: Delivery{
			Name:    "namedel",
			Phone:   "phonedel",
			Zip:     "zepdel",
			City:    "cutydel",
			Address: "adddek",
			Region:  "regdek",
			Email:   "emaldef",
		},
		Payment: Payment{
			Transaction:  "tra",
			RequestId:    "re",
			Currency:     "cu",
			Provider:     "pr",
			Amount:       10,
			PaymentDt:    2745373470,
			Bank:         "bank",
			DeliveryCost: 2,
			GoodsTotal:   3,
			CustomFee:    4,
		},
		Items: []Item{
			{
				ChrtId:      12312,
				TrackNumber: "12321321",
				Price:       231,
				Rid:         "rid",
				Name:        "namememe",
				Sale:        1,
				Size:        "1",
				TotalPrice:  100,
				NmId:        231421312,
				Brand:       "brand",
				Status:      202,
			}, {
				ChrtId:      12312,
				TrackNumber: "12321321",
				Price:       231,
				Rid:         "rid",
				Name:        "namememe",
				Sale:        1,
				Size:        "1",
				TotalPrice:  100,
				NmId:        231421312,
				Brand:       "brand",
				Status:      202,
			},
		},
	}

	data, err := json.Marshal(validOrder)
	if err != nil {
		log.Println(err)
		return
	}

	err = sc.Publish("orders", data)
	if err != nil {
		log.Println(err)
	}
}
