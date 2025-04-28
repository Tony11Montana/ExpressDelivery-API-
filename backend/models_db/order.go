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
	Id_client       uint32  `json:"id_client"`
	Client_name     string  `json:"client_name"`
	Price_delivery  float32 `json:"price_delivery"`
	Type_pay        string  `json:"type_pay"`
	Date_pay        string  `json:"date_pay"`
	Date_order      string  `json:"date_order"`
	Courier_name    string  `json:"courier_name"`
	Count_product   uint64  `json:"count_product"`
	Price_product   float32 `json:"price_product"`
	Name_product    string  `json:"name_product"`
	Count_warehouse uint32  `json:"count_warehouse"`
}

func GetAllOrder(login *string, role *string) ([]*Order, error) {
	var rows *sql.Rows
	var err error

	querydb := `SELECT date_create,date_pay,type_pay,fio,name_product,count_product,price,count_warehouse,fioCourier,priceDelivery
				FROM
				(SELECT orders.*, products.name_product,products.id_warehouse, products.count_warehouse
				FROM 
				(SELECT orders.*, concat(couriers.first_name, " ", couriers.last_name) as fioCourier
				from
				(select orders.*, concat(clients.first_name, " ", clients.last_name) as fio, clients.login
				from
				(select date_create,date_pay,type_pay,id_client as id_client,id_product,count_product,price,id_courier,round((info_orders.price * 0.2), 2) as priceDelivery
				from info_orders inner join orders 
				ON info_orders.id_order = orders.id_order) as orders INNER JOIN clients
				ON orders.id_client = clients.id_client) as orders inner join couriers ON orders.id_courier = couriers.id_courier) as orders INNER JOIN products
				ON orders.id_product = products.id_product ) as orders INNER JOIN warehouses ON orders.id_warehouse = warehouses.id_warehouse`

	if *role == "client" {
		querydb = querydb + ` where login = ?`
		rows, err = db.Query(querydb, &login)
	} else if *role == "courier" {
		rows, err = db.Query(`SELECT date_create,date_pay,type_pay,fio,name_product,count_product,price,count_warehouse,fioCourier,priceDelivery
								FROM
								(SELECT orders.*, products.name_product,products.id_warehouse, products.count_warehouse
								FROM 
								(SELECT orders.*, concat(couriers.first_name, " ", couriers.last_name) as fioCourier
								from
								(select orders.*, concat(clients.first_name, " ", clients.last_name) as fio, clients.login
								from
								(select date_create,date_pay,type_pay,id_client as id_client,id_product,count_product,price,id_courier,round((info_orders.price * 0.2), 2) as priceDelivery
								from info_orders inner join orders 
								ON info_orders.id_order = orders.id_order) as orders INNER JOIN clients
								ON orders.id_client = clients.id_client) as orders inner join couriers ON orders.id_courier = couriers.id_courier where couriers.login = ?) as orders INNER JOIN products
								ON orders.id_product = products.id_product ) as orders INNER JOIN warehouses ON orders.id_warehouse = warehouses.id_warehouse`, &login)
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
			&order.Date_order,
			&order.Date_pay,
			&order.Type_pay,
			&order.Client_name,
			&order.Name_product,
			&order.Count_product,
			&order.Price_product,
			&order.Count_warehouse,
			&order.Courier_name,
			&order.Price_delivery)
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
	var sum_orders float32

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
				sum_orders += product_price
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
				sum_orders += product_price
				countProduct = -countProduct
			}

		}
		fmt.Println(string(product.Product_id), product.Product_name)
	}
	sum_orders *= 0.2 // 20% total sum

	_, err = db.Exec(`update orders set price_delivery = ? where id_order = ?`, &sum_orders, &order.Id_order)
	if err != nil {
		return err
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
