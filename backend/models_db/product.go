package models_db

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Product struct {
	Product_id          uint32  `json:"id"`
	Product_name        string  `json: "product_name"`
	Product_description string  `json: "product_description"`
	Product_price       float32 `json: "product_price"`
	Product_count       uint64  `json: "product_count"`
	Id_warehouse        uint8   `json:"id_warehouse"`
}

func GetAllProducts() ([]*Product, error) {

	//rows, err := db.Query(`select id_product, name_product, description_product, price_product, count_warehouse from products;`)

	/*rows, err := db.Query(`select 0 as id_product, allProducts.name_product, allProducts.description_product,
	round(sum((allProducts.count_warehouse / AllCount.allCount * allProducts.price_product)),2) as price,
	AllCount.allCount
	from (select * from products) as allProducts inner join (select name_product, sum(count_warehouse) as allCount
	from products
	GROUP BY name_product, description_product) as AllCount
	on allProducts.name_product = AllCount.name_product
	group by allProducts.name_product, AllCount.allCount,allProducts.description_product`)*/
	rows, err := db.Query(`SELECT 0,name_product, description_product, price_product, sum(count_warehouse) as count
							FROM products
							group by 
							name_product, description_product, price_product`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	products := make([]*Product, 0)

	var i uint32
	i = 1
	for rows.Next() {
		product := new(Product)
		err := rows.Scan(
			&product.Product_id,
			&product.Product_name,
			&product.Product_description,
			&product.Product_price,
			&product.Product_count)
		if err != nil {
			return nil, err
		}
		product.Product_id = i
		products = append(products, product)
		i++
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}
func AddProduct(pr *Product) error {

	// Check current product.

	//rows, err := db.Query(`select name_product, id_warehouse from Products where name_product = ? AND id_warehouse = ? AND description_product = ?`, &pr.Product_name, &pr.Id_warehouse, &pr.Product_description)

	rows, err := db.Query(`select count_warehouse, price_product from Products where name_product = ? AND id_warehouse = ? AND description_product = ?`, &pr.Product_name, &pr.Id_warehouse, &pr.Product_description)

	if err != nil {
		log.Fatal(err)
		return err
	}

	if !rows.Next() {
		_, err := db.Exec(`insert into products(id_warehouse,name_product,price_product,description_product,count_warehouse) VALUES(?,?,?,?,?)`,
			&pr.Id_warehouse, &pr.Product_name, &pr.Product_price, &pr.Product_description, &pr.Product_count)
		if err != nil {
			log.Fatal(err)
		}
		//return nil
	} else {
		var countProduct uint64
		var priceProduct float32
		//for rows.Next() {
		err := rows.Scan(&countProduct, &priceProduct)
		if err != nil {
			return err
		}
		pr.Product_count += countProduct
		//}
		if err = rows.Err(); err != nil {
			return err
		}

		// average price in warehouse.
		newProductPrice := (priceProduct*float32(countProduct) + (float32(pr.Product_count)-float32(countProduct))*float32(pr.Product_price)) / float32(pr.Product_count)
		//

		fmt.Println(newProductPrice)
		fmt.Println(countProduct)
		fmt.Println(pr.Product_count)

		_, err = db.Exec(`update products set count_warehouse = ?, price_product = ?  where name_product = ? and id_warehouse=? AND description_product = ? `,
			&pr.Product_count, &newProductPrice, &pr.Product_name, &pr.Id_warehouse, &pr.Product_description)
		if err != nil {
			return err
		}
		//return errors.New("This product already exists. You replenished count!")
	}

	var avePrice float32

	rows, err = db.Query(`select round(sum((allProducts.count_warehouse / AllCount.allCount * allProducts.price_product)),2) as price
							from (select * from products) as allProducts inner join (select name_product, sum(count_warehouse) as allCount
							from products 
							GROUP BY name_product, description_product) as AllCount
							on allProducts.name_product = AllCount.name_product
							group by allProducts.name_product, AllCount.allCount,allProducts.description_product
							having name_product = ? and description_product = ?`, &pr.Product_name, &pr.Product_description)
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		err = rows.Scan(&avePrice)
		if err != nil {
			return err
		}

	}
	if err = rows.Err(); err != nil {
		return err
	}
	_, err = db.Exec(`update products set price_product = ? where name_product = ? AND description_product = ? `,
		&avePrice, &pr.Product_name, &pr.Product_description)
	if err != nil {
		return err
	}
	return nil
}
