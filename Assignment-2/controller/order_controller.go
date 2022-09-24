package controller

import (
	"assignment2/constant"
	"assignment2/entity"
	"assignment2/service"
	"assignment2/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type OrderController struct {
	route *mux.Router
	p     *mux.Router
	os    service.OrderService
}

func ProvideController(route *mux.Router, p *mux.Router, os service.OrderService) *OrderController {
	return &OrderController{
		route: route,
		p:     p,
		os:    os,
	}
}

func (o *OrderController) InitController() {
	routes := o.route.PathPrefix(constant.ORDER_API_PATH).Subrouter()
	protected := o.p.Path(constant.ORDER_API_PATH).Subrouter()

	//Order
	protected.HandleFunc("", o.createNewOrder()).Methods(http.MethodPost)
	routes.HandleFunc("", o.viewAllOrders()).Methods(http.MethodGet)
	routes.HandleFunc("/{order_id}", o.deleteOrderByID()).Methods(http.MethodDelete)
	routes.HandleFunc("/{order_id}", o.updateOrderByID()).Methods(http.MethodPut)
}

func (o *OrderController) createNewOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := entity.CreateOrderRequest{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Printf("[createNewOrder] failed to parse JSON data, err => %+v", err)
			panic(err)
		}
		res, err := o.os.CreateNewOrder(r.Context(), &data)
		if err != nil {
			log.Printf("[createNewOrder] failed to create a new order, err => %v", err)
			panic(err)
		}
		utils.NewBaseResponse(http.StatusCreated, "SUCCESS", nil, res).SendResponse(&w)
	}
}

func (o *OrderController) viewAllOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := o.os.ViewAllOrders(r.Context())
		if err != nil {
			log.Printf("[viewAllOrders] failed to get all the orders, err => %v", err)
			panic(err)
		}
		utils.NewBaseResponse(http.StatusOK, "SUCCESS", nil, res).SendResponse(&w)
	}
}
func (o *OrderController) deleteOrderByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routeVar := mux.Vars(r)
		idVar := routeVar["order_id"]
		idConv, _ := strconv.ParseUint(idVar, 10, 64)

		res, err := o.os.DeleteOrderByID(r.Context(), idConv)
		if err != nil {
			log.Printf("[deleteOrderByID] failed to delete the order by id, err => %v, id => %v", err, idConv)
			panic(err)
		}
		utils.NewBaseResponse(http.StatusAccepted, "SUCCESS", nil, res).SendResponse(&w)
	}
}

func (o *OrderController) updateOrderByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routeVar := mux.Vars(r)
		idVar := routeVar["order_id"]
		idConv, _ := strconv.ParseUint(idVar, 10, 64)

		data := new(entity.UpdateOrderRequest)
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Printf("[updateOrderByID] failed to parse JSON data, err => %+v", err)
			panic(err)
		}

		res, err := o.os.UpdateOrderByID(r.Context(), idConv, data)
		if err != nil {
			log.Printf("[updateOrderByID] failed to update the order by id, err => %v, id => %v", err, idConv)
			panic(err)
		}
		utils.NewBaseResponse(http.StatusCreated, "SUCCESS", nil, res).SendResponse(&w)
	}
}
