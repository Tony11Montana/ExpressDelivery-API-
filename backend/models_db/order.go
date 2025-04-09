package models_db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type OrderCreate struct {
	id_client         uint32
	date_order        string
	type_pay          string
	date_pay          string
	description_order string
	price_delivery    uint64
	Products_in_order []Product 
}

type Order struct {
	Id_client       uint32 `json:"id_client"`
	Client_name     string `json:"client_name"`
	Price_delivery  uint32 `json:"price_delivery"`
	Type_pay        string `json:"type_pay"`
	Date_pay        string `json:"date_pay"`
	Date_order      string `json:"date_order"`
	Courier_name    string `json:"courier_name"`
	Count_product   uint64 `json:"count_product"`
	Price_product   uint64 `json:"price_product"`
	Name_product    string `json:"name_product"`
	Count_warehouse uint32 `json:"count_warehouse"`
}

func GetAllOrder(login *string, role *string) ([]*Order, error) {
	var rows *sql.Rows
	var err error

	querydb := `SELECT Orders_new.id_client, concat(Clients.first_name," " ,Clients.last_name) as client_name,
						Orders_new.price_delivery, Orders_new.type_pay, Orders_new.date_pay, Orders_new.date_order,
						Orders_new.CourierName as courier_name, Orders_new.count_product, Orders_new.price as price_product,
						Orders_new.name_product, Orders_new.count_warehouse
						from (
						select Orders.id_client, Orders.price_delivery, Orders.type_pay, Orders.date_pay, Orders.date_order,
						info_cour.CourierName, info_cour.count_product, info_cour.price,
						info_cour.name_product, info_cour.count_warehouse
						from Orders inner join (
						select concat(Couriers.first_name, " " , Couriers.last_name) as CourierName,
						Info_orders_new.count_product, Info_orders_new.price, Info_orders_new.id_order,
						Info_orders_new.name_product, Info_orders_new.count_warehouse
						from (select Info_orders.id_courier, Info_orders.count_product, Info_orders.price, Info_orders.id_order,
						Products.name_product, Products.count_warehouse
						from Info_orders INNER join Products
						) as Info_orders_new INNER JOIN Couriers ON Info_orders_new.id_courier = Couriers.id_courier
						) as info_cour ON Orders.id_order = info_cour.id_order
						) as Orders_new INNER JOIN Clients ON
						Orders_new.id_client = Clients.id_client`

	if *role == "client" {
		querydb = querydb + ` where clients.login = ?`
		rows, err = db.Query(querydb, &login)
	} else {
		rows, err = db.Query(querydb)
	}

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer rows.Close()

	orders := make([]*Order, 0)
	for rows.Next() {
		order := new(Order)
		err := rows.Scan(
			&order.Id_client,
			&order.Client_name,
			&order.Price_delivery,
			&order.Type_pay,
			&order.Date_pay,
			&order.Date_order,
			&order.Courier_name,
			&order.Count_product,
			&order.Price_product,
			&order.Name_product,
			&order.Count_warehouse)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

/*
	rows, err := db.DB.Query(`SELECT * from ElaginDiplom.Products`)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		fmt.Println(rows)
	}

func (p *Order) AddOrder() {

}

func (p *Order) UpdateOrder() {

}*/
