package entity

import "time"

// order_request_dto
type Orders struct {
	OrderID      uint64    `db:"order_id"`
	CustomerName string    `db:"customer_name"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Items struct {
	ItemId      uint64    `db:"item_id"`
	ItemCode    string    `db:"item_code"`
	Description string    `db:"description"`
	Quantity    uint64    `db:"quantity"`
	OrderID     uint64    `db:"order_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type OrdersItemsJoined struct {
	Orders
	Items AllItems
}

type OrdersJoined []*OrdersItemsJoined

type OrderRequest struct {
	OrderID int `json:"order_id"`
}

type AllOrders []*Orders
type AllItems []*Items

type UpdateOrderRequest struct {
	CustomerName string `json:"customer_name"`
	ItemsRequest `json:"items"`
}

type CreateOrderRequest struct {
	CustomerName string `json:"customer_name"`
	ItemsRequest `json:"items"`
}

type ItemRequest struct {
	ItemID      uint64 `json:"item_id"`
	ItemCode    string `json:"item_code"`
	Description string `json:"description"`
	Quantity    uint64 `json:"quantity"`
}

type ItemsRequest []ItemRequest

type UpdateOrdersByIDResponse struct {
	CustomerName string    `json:"customer_name"`
	UpdatedAt    time.Time `json:"updated_at"`
	ItemsUpdate  `json:"items"`
}

type ItemUpdateResponse struct {
	ItemID      uint64 `json:"item_id"`
	ItemCode    string `json:"item_code"`
	Description string `json:"description"`
	Quantity    uint64 `json:"quantity"`
}

type OrderResponse struct {
	OrderID          uint64    `json:"order_id"`
	CustomerName     string    `json:"customer_name"`
	CreatedAt        time.Time `json:"ordered_at"`
	AllItemsResponse `json:"items"`
}

type ItemsResponse struct {
	ItemID      uint64 `json:"item_id"`
	ItemCode    string `json:"item_code"`
	Description string `json:"description"`
	Quantity    uint64 `json:"quantity"`
	OrderID     uint64 `json:"order_id"`
}

type OrderDetails []OrderResponse
type AllItemsResponse []ItemsResponse
type ItemsUpdate []ItemUpdateResponse

func CreateItemUpdateResponse(i Items) ItemUpdateResponse {
	return ItemUpdateResponse{
		ItemID:      i.ItemId,
		ItemCode:    i.ItemCode,
		Description: i.Description,
		Quantity:    i.Quantity,
	}
}

func CreateOrderResponse(order Orders) OrderResponse {
	return OrderResponse{
		CustomerName: order.CustomerName,
		CreatedAt:    time.Now(),
	}
}

func CreateItemsResponse(item Items) ItemsResponse {
	return ItemsResponse{
		ItemID:      item.ItemId,
		ItemCode:    item.ItemCode,
		Description: item.Description,
		Quantity:    item.Quantity,
	}
}

func CreateUpdateResponse(o Orders, is AllItems) *UpdateOrdersByIDResponse {
	updateDetails := UpdateOrdersByIDResponse{}
	updateDetails.CustomerName = o.CustomerName
	updateDetails.UpdatedAt = time.Now()
	for idx := range is {
		allItems := CreateItemUpdateResponse(*is[idx])
		updateDetails.ItemsUpdate = append(updateDetails.ItemsUpdate, allItems)
	}
	return &updateDetails
}

func CreateOrderResponseDetail(o Orders, is AllItems, id uint64) *OrderResponse {
	orderDetails := CreateOrderResponse(o)
	orderDetails.OrderID = id
	for idx := range is {
		items := CreateItemsResponse(*is[idx])
		items.OrderID = id
		orderDetails.AllItemsResponse = append(orderDetails.AllItemsResponse, items)
	}

	return &orderDetails
}

func ViewOrderResponseDetails(os OrdersJoined) []OrderResponse {
	var orderDetails []OrderResponse

	for _, each := range os {
		order := viewOrderResponse(*each)
		orderDetails = append(orderDetails, order)
	}
	return orderDetails
}

func viewOrderResponse(os OrdersItemsJoined) OrderResponse {
	var listItems AllItemsResponse

	for _, each := range os.Items {
		items := CreateItemsResponse(*each)
		listItems = append(listItems, items)
	}

	return OrderResponse{
		OrderID:          os.OrderID,
		CustomerName:     os.CustomerName,
		CreatedAt:        os.CreatedAt,
		AllItemsResponse: listItems,
	}
}

// JSON --> entity
func (or *CreateOrderRequest) ToEntity() (item AllItems, order Orders) {
	item = AllItems{}
	for _, items := range or.ItemsRequest {
		itemsDetail := Items{
			ItemCode:    items.ItemCode,
			Description: items.Description,
			Quantity:    items.Quantity,
		}
		item = append(item, &itemsDetail)
	}

	order = Orders{
		CustomerName: or.CustomerName,
	}
	return
}

func (or *UpdateOrderRequest) ToEntity() (order Orders, item AllItems) {
	order = Orders{
		CustomerName: or.CustomerName,
	}

	item = AllItems{}
	for _, items := range or.ItemsRequest {
		itemsDetails := Items{
			ItemId:      items.ItemID,
			ItemCode:    items.ItemCode,
			Description: items.Description,
			Quantity:    items.Quantity,
		}
		item = append(item, &itemsDetails)
	}

	return
}
