package service

import (
	"assignment2/entity"
	"assignment2/repository"
	"context"
	"log"
)

type OrderServiceImpl struct {
	repo repository.OrderRepository
}

type OrderService interface {
	CreateNewOrder(ctx context.Context, data *entity.CreateOrderRequest) (*entity.OrderResponse, error)
	ViewAllOrders(ctx context.Context) (*entity.OrderDetails, error)
	DeleteOrderByID(ctx context.Context, id uint64) (int, error)
	UpdateOrderByID(ctx context.Context, id uint64, data *entity.UpdateOrderRequest) (*entity.UpdateOrdersByIDResponse, error)
}

func ProvideService(repo repository.OrderRepository) *OrderServiceImpl {
	return &OrderServiceImpl{
		repo: repo,
	}
}

func (o OrderServiceImpl) UpdateOrderByID(ctx context.Context, id uint64, data *entity.UpdateOrderRequest) (*entity.UpdateOrdersByIDResponse, error) {
	check, err := o.repo.CheckOrders(ctx)
	if err != nil {
		log.Printf("[UpdateOrderByID] an error occured while checking orders, err => %v, id => %v", err, id)
		return nil, err
	}
	if check < 1 {
		log.Printf("[DeleteOrderByID] theres is no orders data, err => %v", err)
		panic(err)
	}

	order, item := data.ToEntity()
	err = o.repo.UpdateOrderByID(ctx, id, &order, item)
	if err != nil {
		log.Printf("[UpdateOrderByID] an error occured while updating orders, err => %v, id => %v", err, id)
		return nil, err
	}

	response := entity.CreateUpdateResponse(order, item)
	return response, nil
}

func (o OrderServiceImpl) DeleteOrderByID(ctx context.Context, id uint64) (int, error) {
	check, err := o.repo.CheckOrders(ctx)
	if err != nil {
		log.Printf("[DeleteOrderByID] an error occured while checking orders, err => %v, id => %v", err, id)
		return 0, nil
	}
	if check < 1 {
		log.Printf("[DeleteOrderByID] theres is no orders data, err => %v", err)
		panic(err)
	}
	res, err := o.repo.DeleteOrderByID(ctx, id)
	if err != nil {
		log.Printf("[DeleteOrderByID] an error occured while deleting orders, err => %v, id => %v", err, id)
		return 0, nil
	}

	return res, nil
}

func (o OrderServiceImpl) CreateNewOrder(ctx context.Context, data *entity.CreateOrderRequest) (*entity.OrderResponse, error) {
	item, order := data.ToEntity()
	orderID, err := o.repo.CreateNewOrder(ctx, order, item)
	if err != nil {
		log.Printf("[CreateNewOrder] an error occured while creating new order, err => %v", err)
		return nil, err
	}

	return entity.CreateOrderResponseDetail(order, item, orderID), nil
}

func (o OrderServiceImpl) ViewAllOrders(ctx context.Context) (*entity.OrderDetails, error) {
	count, err := o.repo.CheckOrders(ctx)
	if err != nil {
		log.Printf("[ViewAllOrders] an error occured while show all the orders, err => %v", err)
		return nil, err
	}
	if count < 1 {
		log.Printf("[ViewAllOrders] theres is no orders data, err => %v", err)
		panic(err)

	}
	res, err := o.repo.ViewAllOrders(ctx)
	if err != nil {
		log.Printf("[ViewAllOrders] an error occured while show all the orders, err => %v", err)
		return nil, err
	}

	var response1 entity.OrderDetails
	response := entity.ViewOrderResponseDetails(res)
	response1 = append(response1, response...)

	return &response1, nil
}
