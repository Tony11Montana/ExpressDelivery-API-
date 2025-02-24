package models_db

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Courier struct {
	Id_courier     uint16 `json:"id_courier"`
	First_name     string `json:"first_name"`
	Last_name      string `json:"last_name"`
	Total_salary   uint32 `json:"total_salary"`
	Warehouse_name string `json:"warehouse_name"`
	Id_warehouse   uint8  `json:"id_warehouse"`
}

func GetAllCouriers() ([]*Courier, error) {

	rows, err := db.Query(`SELECT subQueryCourier.id_courier, first_name, last_name, ifnull(sum(price_delivery), 0) as total_salary,
	                       ifnull(name_warehouse, "-") as name_warehouse, ifnull(subQueryCourier.id_warehouse, 0) as id_warehouse				   
						   from (
						   Select Couriers.id_courier, Couriers.first_name, Couriers.last_name, warehouses.name_warehouse, warehouses.id_warehouse 
						   from couriers left join warehouses on
						   couriers.id_warehouse = warehouses.id_warehouse
						   ) as subQueryCourier left JOIN (
						   SELECT Orders.price_delivery, Info_orders.id_order, Info_orders.id_courier  
						   FROM Orders INNER JOIN Info_orders
						   ON Orders.id_order = Info_orders.id_order) as subQuery
						   ON subQueryCourier.id_courier = subQuery.id_courier
						   GROUP BY 
						   first_name, 
						   last_name,
						   Couriers.id_courier,
						   name_warehouse,
                           id_warehouse`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	couriers := make([]*Courier, 0)

	for rows.Next() {
		courier := new(Courier)
		err := rows.Scan(
			&courier.Id_courier,
			&courier.First_name,
			&courier.Last_name,
			&courier.Total_salary,
			&courier.Warehouse_name,
			&courier.Id_warehouse)
		if err != nil {
			return nil, err
		}
		couriers = append(couriers, courier)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return couriers, nil
}

func AddCourier(courier *Courier) (err error) {
	_, err = db.Exec(`insert into Couriers(first_name,last_name,id_warehouse) values(?, ?, ?)`, &courier.First_name, &courier.Last_name, &courier.Id_warehouse)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
