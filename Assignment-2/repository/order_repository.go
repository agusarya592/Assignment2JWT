package repository

import (
	"assignment2/entity"
	"context"
	"database/sql"
	"log"
)

type OrderRepositoryImpl struct {
	DB *sql.DB
}


type OrderRepository interface {
	CreateNewOrder(ctx context.Context, reqDataOrder entity.Orders, reqDataItems entity.AllItems) (uint64, error)
	ViewAllOrders(ctx context.Context) (entity.OrdersJoined, error)
	CheckOrders(ctx context.Context) (int, error)
	DeleteOrderByID(ctx context.Context, id uint64) (int, error)
	UpdateOrderByID(ctx context.Context, id uint64, reqDataOrder *entity.Orders, reqDataItems entity.AllItems) error
	GetOrdersByID(ctx context.Context, id uint64) (entity.OrdersItemsJoined, error)
}

func ProvideRepository(DB *sql.DB) *OrderRepositoryImpl {
	return &OrderRepositoryImpl{DB: DB}
}

var (
	INSERT_ORDER_DATA   = "INSERT INTO `order` (customer_name) VALUES (?)"
	INSERT_ITEM_DATA    = "INSERT INTO `item` (item_code, description, quantity, order_id) VALUES(?, ?, ?, ?)"
	SELECT_ORDERS       = "SELECT o.order_id, o.customer_name, o.created_at, o.updated_at FROM `order` o"
	SELECT_ORDERS_BY_ID = "SELECT o.order_id, o.customer_name, o.created_at, o.updated_at FROM `order` o WHERE o.order_id = ?"
	SELECT_ITEMS        = "SELECT i.item_id, i.item_code, i.description, i.quantity, i.order_id FROM `item` i WHERE i.order_id=?"
	COUNT_ORDERS        = "SELECT COUNT(*) FROM `order`"
	COUNT_ITEMS         = "SELECT COUNT(*) FROM `item` i WHERE i.order_id = ?"
	DELETE_ORDER        = "DELETE FROM `order` WHERE order_id = ?"
	DELETE_ITEMS        = "DELETE FROM `item` WHERE order_id = ?"
	UPDATE_ORDER        = "UPDATE `order` SET customer_name = ? WHERE order_id = ?"
	UPDATE_ITEM         = "UPDATE item SET description = ?, item_code = ?, quantity = ? WHERE item_id = ?"
)



func (o OrderRepositoryImpl) GetOrdersByID(ctx context.Context, id uint64) (entity.OrdersItemsJoined, error) {
	query := SELECT_ORDERS_BY_ID

	tx, err := o.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[GetOrdersByID] failed to begin transaction, err => %v", err)
		return entity.OrdersItemsJoined{}, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("[GetOrdersByID] failed to prepare the statement, err => %v", err)
		return entity.OrdersItemsJoined{}, err
	}

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		log.Printf("[GetOrdersByID] failed to query to the database, err => %v", err)
		return entity.OrdersItemsJoined{}, err
	}

	var orders entity.OrdersItemsJoined
	for rows.Next() {
		orderDetails := entity.OrdersItemsJoined{}

		err := rows.Scan(
			&orderDetails.Orders.OrderID,
			&orderDetails.Orders.CustomerName,
			&orderDetails.Orders.CreatedAt,
			&orderDetails.Orders.UpdatedAt,
		)
		if err != nil {
			log.Printf("[GetOrdersByID] failed to scan the data, err => %v", err)
			return entity.OrdersItemsJoined{}, err
		}
		orders = orderDetails
	}
	countItems := COUNT_ITEMS

	stmt, err = o.DB.PrepareContext(ctx, countItems)
	if err != nil {
		log.Printf("[GetOrdersByID] failed to prepare the statement, err => %v", err)
		return entity.OrdersItemsJoined{}, err
	}

	rows, err = stmt.QueryContext(ctx, id)
	if err != nil {
		log.Printf("[GetOrdersByID] failed to query to the database, err => %v", err)
		return entity.OrdersItemsJoined{}, err
	}

	var itemCount int

	for rows.Next() {
		err := rows.Scan(
			&itemCount,
		)
		if err != nil {
			log.Printf("[GetOrdersByID] failed to scan the data, err => %v", err)
			return entity.OrdersItemsJoined{}, err
		}
	}
	queryItems := SELECT_ITEMS
	for idx := 1; idx < itemCount; idx++ {
		stmt, err := tx.PrepareContext(ctx, queryItems)
		if err != nil {
			log.Printf("[GetOrdersByID] failed to prepare the statement, err => %v", err)
			return entity.OrdersItemsJoined{}, err
		}

		rows, err := stmt.QueryContext(ctx, orders.OrderID)
		if err != nil {
			log.Printf("[GetOrdersByID] failed to query to the database, err => %v", err)
			return entity.OrdersItemsJoined{}, err
		}

		for rows.Next() {
			allItems := entity.Items{}

			err := rows.Scan(
				&allItems.ItemId,
				&allItems.ItemCode,
				&allItems.Description,
				&allItems.Quantity,
				&allItems.OrderID,
			)
			if err != nil {
				log.Printf("[GetOrdersByID] failed to scan the data, err => %v", err)
				return entity.OrdersItemsJoined{}, err
			}
			orders.Items = append(orders.Items, &allItems)
		}
	}
	return orders, nil
}

func (o OrderRepositoryImpl) CreateNewOrder(ctx context.Context, reqDataOrder entity.Orders, reqDataItems entity.AllItems) (uint64, error) {
	tx, err := o.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[CreateNewOrder] failed to begin transaction, err => %v", err)
		return 0, err
	}
	defer tx.Rollback()

	queryOrder := INSERT_ORDER_DATA
	res, err := tx.ExecContext(ctx, queryOrder, reqDataOrder.CustomerName)
	if err != nil {
		log.Printf("[CreateNewOrder] failed to insert order data, err => %v", err)
		return 0, err
	}

	orderID, _ := res.LastInsertId()

	for _, items := range reqDataItems {
		queryItemData := INSERT_ITEM_DATA
		stmt, err := tx.PrepareContext(ctx, queryItemData)
		if err != nil {
			log.Printf("[CreateNewOrder] failed to prepare the statement, err => %v", err)
			return 0, err
		}
		res, err = stmt.ExecContext(
			ctx,
			items.ItemCode,
			items.Description,
			items.Quantity,
			orderID,
		)
		if err != nil {
			log.Printf("[CreateNewOrder] failed to insert items data, err => %v", err)
			return 0, err
		}

		itemsID, _ := res.LastInsertId()
		items.ItemId = uint64(itemsID)
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[CreateNewOrder] transaction failed, err => %v", err)
		return 0, err
	}
	return uint64(orderID), nil
}

func (o OrderRepositoryImpl) UpdateOrderByID(ctx context.Context, id uint64, reqDataOrder *entity.Orders, reqDataItems entity.AllItems) error {
	queryItems := UPDATE_ITEM
	queryOrder := UPDATE_ORDER

	stmt, err := o.DB.PrepareContext(ctx, queryOrder)
	if err != nil {
		log.Printf("[UpdateOrderByID] failed to prepare the statement, err => %v", err)
		return err
	}

	_, err = stmt.ExecContext(
		ctx,
		reqDataOrder.CustomerName,
		id,
	)
	if err != nil {
		log.Printf("[UpdateOrderByID] failed to update the data, err => %v", err)
		return err
	}

	for _, items := range reqDataItems {
		stmt, err = o.DB.PrepareContext(ctx, queryItems)
		if err != nil {
			log.Printf("[UpdateOrderByID] failed to prepare the statement, err => %v", err)
			return err
		}
		_, err = stmt.ExecContext(
			ctx,
			items.Description,
			items.ItemCode,
			items.Quantity,
			items.ItemId,
		)
		if err != nil {
			log.Printf("[UpdateOrderByID] failed to update the data, err => %v", err)
			return err
		}
	}
	return nil
}

func (o OrderRepositoryImpl) DeleteOrderByID(ctx context.Context, id uint64) (int, error) {
	tx, err := o.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[DeleteOrderByID] failed to begin transaction, err => %v", err)
		return 0, err
	}
	defer tx.Rollback()

	queryDeleteItem := DELETE_ITEMS
	_, err = o.DB.Query(queryDeleteItem, id)
	if err != nil {
		log.Printf("[DeleteByOrderID] failed to delete item data, err => %v", err)
		return 0, err
	}

	queryDeleteOrder := DELETE_ORDER
	_, err = o.DB.Query(queryDeleteOrder, id)
	if err != nil {
		log.Printf("[DeleteOrderByID] failed to delete order data, err => %v", err)
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[DeleteOrderByID] transaction failed, err => %v", err)
		return 0, err
	}
	return int(id), nil
}

func (o OrderRepositoryImpl) CheckOrders(ctx context.Context) (int, error) {
	query := COUNT_ORDERS

	stmt, err := o.DB.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("[CheckOrders] failed to prepare the statement, err => %v", err)
		return 0, err
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		log.Printf("[CheckOrders] failed to query to the database, err => %v", err)
		return 0, err
	}

	var orderCount int

	for rows.Next() {
		err := rows.Scan(
			&orderCount,
		)
		if err != nil {
			log.Printf("[CheckOrders] failed to scan data from database, err => %v", err)
			return 0, err
		}
	}
	return orderCount, nil
}

func (o OrderRepositoryImpl) ViewAllOrders(ctx context.Context) (entity.OrdersJoined, error) {
	query := SELECT_ORDERS

	stmt, err := o.DB.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("[ViewAllOrders] failed to prepare the statement, err => %v", err)
		return nil, err
	}
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		log.Printf("[ViewAllOrders] failed to query to the database, err => %v", err)
		return nil, err
	}

	var ordersJoined entity.OrdersJoined

	for rows.Next() {
		orderItemJoined := entity.OrdersItemsJoined{}

		err := rows.Scan(
			&orderItemJoined.Orders.OrderID,
			&orderItemJoined.Orders.CustomerName,
			&orderItemJoined.Orders.CreatedAt,
			&orderItemJoined.Orders.UpdatedAt,
		)
		if err != nil {
			log.Printf("[ViewAllOrders] failed to scan data from database, err => %v", err)
			return nil, err
		}

		ordersJoined = append(ordersJoined, &orderItemJoined)
	}
	query = SELECT_ITEMS
	for idx, items := range ordersJoined {

		stmt, err = o.DB.PrepareContext(ctx, query)
		if err != nil {
			log.Printf("[ViewAllOrders] failed to prepare the statement, err => %v", err)
			return nil, err
		}
		rows, err := stmt.QueryContext(ctx, ordersJoined[idx].Orders.OrderID)
		if err != nil {
			log.Printf("[ViewAllOrders] failed to query to the database, err => %v", err)
			return nil, err
		}

		for rows.Next() {
			allItems := entity.Items{}

			err = rows.Scan(
				&allItems.ItemId,
				&allItems.ItemCode,
				&allItems.Description,
				&allItems.Quantity,
				&allItems.OrderID,
			)
			if err != nil {
				log.Printf("[ViewAllOrders] failed to scan data from database, err => %v", err)
				return nil, err
			}

			items.Items = append(items.Items, &allItems)
		}
	}
	return ordersJoined, nil
}
