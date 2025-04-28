package models_db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Courier struct {
	Id_courier     uint16  `json:"id_courier"`
	First_name     string  `json:"first_name"`
	Last_name      string  `json:"last_name"`
	Total_salary   float32 `json:"total_salary"`
	Warehouse_name string  `json:"warehouse_name"`
	Id_warehouse   uint8   `json:"id_warehouse"`
}

func GetAllCouriers(login *string, role *string) ([]*Courier, error) {
	var rows *sql.Rows
	var err error

	if *role == "courier" {
		rows, err = db.Query(`select id_courier,first_name,last_name,sum(total_sum),warehouses.name_warehouse, warehouses.id_warehouse
							from
							(select couriers.id_courier,couriers.first_name, couriers.last_name, ifnull(round(info_orders.price * 0.2, 2), 0) as total_sum, couriers.id_warehouse
							from info_orders right join couriers 
							on info_orders.id_courier = couriers.id_courier where login = ?) as couriers inner join warehouses on couriers.id_warehouse = warehouses.id_warehouse
							group by
							id_courier,
							first_name,
							last_name,
							id_warehouse,
							name_warehouse`, &login)
	} else {
		rows, err = db.Query(`select id_courier,first_name,last_name,sum(total_sum),warehouses.name_warehouse, warehouses.id_warehouse
							from
							(select couriers.id_courier,couriers.first_name, couriers.last_name, ifnull(round(info_orders.price * 0.2, 2), 0) as total_sum, couriers.id_warehouse
							from info_orders right join couriers 
							on info_orders.id_courier = couriers.id_courier) as couriers inner join warehouses on couriers.id_warehouse = warehouses.id_warehouse
							group by
							id_courier,
							first_name,
							last_name,
							id_warehouse,
							name_warehouse`)
	}
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
