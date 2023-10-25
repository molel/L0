package repo

import (
	"L0/internal/entity"
	"encoding/json"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"sync"
)

const (
	getAllOrders = "SELECT * from orders AS o"
	getOrderById = "SELECT o.order FROM orders AS o WHERE order_uid = $1;"
	insertOrder  = "INSERT INTO orders VALUES ($1, $2);"
)

type ordersCache struct {
	sync.RWMutex
	m map[string]entity.UnmarshalledOrder
}

func newOrdersCache() *ordersCache {
	return &ordersCache{m: make(map[string]entity.UnmarshalledOrder)}
}

func (oc *ordersCache) Set(key string, value entity.UnmarshalledOrder) {
	oc.Lock()
	defer oc.Unlock()
	oc.m[key] = value
}

func (oc *ordersCache) Get(key string) (entity.UnmarshalledOrder, bool) {
	oc.RLock()
	defer oc.RUnlock()
	value, ok := oc.m[key]
	return value, ok
}

type OrderRepo struct {
	db    *sqlx.DB
	cache *ordersCache
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	repo := &OrderRepo{db: db, cache: newOrdersCache()}

	err := repo.restoreCache()
	if err != nil {
		log.Fatalf("Error occurred during restoring cache from db: %s", err)
	}

	return repo
}

func (repo *OrderRepo) restoreCache() error {
	var rawOrders []entity.RawOrder
	err := repo.db.Select(&rawOrders, getAllOrders)

	if err != nil {
		return err
	}

	var unmarshalledOrder entity.UnmarshalledOrder
	for _, rawOrder := range rawOrders {
		if err := json.Unmarshal(rawOrder.Order, &unmarshalledOrder); err != nil {
			repo.cache.Set(unmarshalledOrder.OrderUid, unmarshalledOrder)
		}
	}

	return nil
}

func (repo *OrderRepo) GetOrderById(orderUid string) (entity.UnmarshalledOrder, error) {
	if order, ok := repo.cache.Get(orderUid); ok {
		return order, nil
	}

	var rawOrder []byte
	if err := repo.db.Get(&rawOrder, getOrderById, orderUid); err != nil {
		return entity.UnmarshalledOrder{}, err
	}

	var unmarshalledOrder entity.UnmarshalledOrder
	if err := json.Unmarshal(rawOrder, &unmarshalledOrder); err != nil {
		return entity.UnmarshalledOrder{}, err
	}

	if _, ok := repo.cache.Get(orderUid); !ok {
		repo.cache.Set(orderUid, unmarshalledOrder)
	}

	return unmarshalledOrder, nil
}

func (repo *OrderRepo) InsertOrder(orderUid string, unmarshalledOrder entity.UnmarshalledOrder) error {
	rawOrder, err := json.Marshal(unmarshalledOrder)
	if err != nil {
		log.Printf("Cannot marshall order: %s", err)
	}

	_, err = repo.db.Exec(insertOrder, orderUid, rawOrder)

	if err == nil {
		var unmarshalledOrder entity.UnmarshalledOrder
		if err = json.Unmarshal(rawOrder, &unmarshalledOrder); err != nil {
			repo.cache.Set(orderUid, unmarshalledOrder)
		}
	}

	return err
}
