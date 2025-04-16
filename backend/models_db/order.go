package models_db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type OrderCreate struct {
	Id_order          uint8
	Id_client         uint8
	Id_employee       uint8
	Date_order        string `json:"creationDate"`
	Type_pay          string `json:"paymentType"`
	Date_pay          string `json:"paymentDate"`
	Description_order string `json:"orderDescription"`
	Price_delivery    uint64
	Products_in_order []Product `json:"products"`
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

func CreateOrder(login *string, order *OrderCreate) error {
	//get client id.
	rows, err := db.Query(`select id_client from clients where login = ?`, &login)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&order.Id_client)
		if err != nil {
			return err
		}
	}
	if err = rows.Err(); err != nil {
		return err
	}

	//get id employee.
	rows, err = db.Query(`select employee.id_employee from (select orders.id_employee as id_employee, count(orders.id_employee) as efficiency, employees.first_name as first_name from orders right join employees on orders.id_employee=employees.id_employee group by id_employee, employees.first_name having id_employee is not NULL order by efficiency ASC LIMIT 1) as employee`)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&order.Id_employee)
		if err != nil {
			return err
		}
	}
	if err = rows.Err(); err != nil {
		return err
	}

	count_warehouse := 0

	//Count warehouses
	rows, err = db.Query(`select count(*) as count_warehouses from warehouses`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&count_warehouse)
		if err != nil {
			return err
		}
	}
	if err = rows.Err(); err != nil {
		return err
	}

	rows, err = db.Query(`select id_order + 1 from orders order by id_order DESC limit 1`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&order.Id_order)
		if err != nil {
			return err
		}
	}
	if err = rows.Err(); err != nil {
		return err
	}

	_, err = db.Exec(`insert into orders (id_order,id_client,id_employee,id_status,date_order,type_pay,date_pay,description_order,price_delivery ) values (?,?,?,?,?,?,?,?,?)`,
		&order.Id_order, &order.Id_client, &order.Id_employee, 1, &order.Date_order, &order.Type_pay, &order.Date_pay, &order.Description_order, 0)
	if err != nil {
		return err
	}

	fmt.Println("OrderPage", order.Products_in_order[0].Product_count)
	//seach courier of warehouse and few work. ( for each item ) and create info_orders
	for _, product := range order.Products_in_order {
		countProductAll := product.Product_count
		countProduct := countProductAll
		for i := 0; i < count_warehouse; i++ {
			id_warehouse := 0
			id_courier := 0
			id_product := 0
			rows, err = db.Query(`select products.id_product, products.id_warehouse, products.count_warehouse 
									from products inner join warehouses 
									on products.id_warehouse = warehouses.id_warehouse 
									where products.name_product = ?
									order by products.count_warehouse DESC limit 1`, &product.Product_name)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				err = rows.Scan(
					&id_product,
					&id_warehouse,
					&countProductAll)
				if err != nil {
					return err
				}
			}
			if err = rows.Err(); err != nil {
				return err
			}
			//search courier of few work.
			rows, err := db.Query(`select couriers.id_courier
									from (select id_courier as idCourier, count(id_courier) as listwork
									from info_orders
									group by id_courier) as couriers_fewWork right join couriers
									on couriers_fewWork.idCourier = couriers.id_courier
									where id_warehouse = ?
									order by couriers_fewWork.listwork asc
									limit 1`, &id_warehouse)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				err = rows.Scan(
					&id_courier)
				if err != nil {
					return err
				}
			}
			if err = rows.Err(); err != nil {
				return err
			}

			countProduct = countProductAll - countProduct
			if countProduct >= 0 {
				_, err = db.Exec(`update products set count_warehouse = ? where name_product = ? and id_warehouse=?`,
					&countProduct, &product.Product_name, &id_warehouse)
				if err != nil {
					return err
				}
				product_price := product.Product_price * float32(countProductAll-countProduct)
				count := countProductAll - countProduct
				_, err = db.Exec(`insert into info_orders(id_courier,id_order,id_product,count_product,price,date_create ) VALUES(?,?,?,?,?,?)`,
					&id_courier, &order.Id_order, &id_product, &count, &product_price, &order.Date_order)
				if err != nil {
					return err
				}
				break
			} else if countProduct < 0 {
				_, err = db.Exec(`update products set count_warehouse = ? where name_product = ? and id_warehouse=?`,
					0, &product.Product_name, &id_warehouse)
				if err != nil {
					return err
				}
				product_price := product.Product_price * float32(countProductAll)
				_, err = db.Exec(`insert into info_orders(id_courier,id_order,id_product,count_product,price,date_create ) VALUES(?,?,?,?,?,?)`,
					&id_courier, &order.Id_order, &id_product, &countProductAll, &product_price, &order.Date_order)
				if err != nil {
					return err
				}
				countProduct = -countProduct
			}

		}
		fmt.Println(string(product.Product_id), product.Product_name)
	}

	return nil

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
