package app

import (
	"L0/config"
	httpController "L0/internal/controller/http"
	natsController "L0/internal/controller/nats-streaming"
	"L0/internal/repo"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func Run(cfg config.Config) {

	dbConnURL := fmt.Sprintf(
		"postgresql://%s:%s@localhost:5432/%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDb,
	)

	db, err := sqlx.Connect("pgx", dbConnURL)
	if err != nil {
		log.Fatalf("Error occurred during connecting to database: %s", err)
	}

	repository := repo.NewOrderRepo(db)

	controller := httpController.NewOrderController(repository)

	subscriber := natsController.NewOrderSubscriber(repository)
	subscriber.Subscribe(cfg.StanClusterId)

	http.HandleFunc("/order", controller.GetOrder)

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", cfg.HttpPort), nil))
}
