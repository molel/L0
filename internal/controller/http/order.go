package http

import (
	"L0/internal/repo"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

var (
	orderTemplate = template.Must(template.ParseFiles(".\\..\\web\\views\\order.html"))
)

type OrderController struct {
	repo *repo.OrderRepo
}

func NewOrderController(repo *repo.OrderRepo) *OrderController {
	return &OrderController{repo: repo}
}

func (controller *OrderController) GetOrder(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	orderUid := request.URL.Query().Get("order_uid")

	unmarshalledOrder, err := controller.repo.GetOrderById(orderUid)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		http.Error(writer, fmt.Sprintf("Error occurred during getting rawOrder from database: %s", err), http.StatusInternalServerError)
		return
	}

	if err = orderTemplate.Execute(writer, unmarshalledOrder); err != nil {
		http.Error(writer, fmt.Sprintf("Error occurred during executing the template: %s", err), http.StatusInternalServerError)
		return
	}
}
