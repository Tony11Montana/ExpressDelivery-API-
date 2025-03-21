package models_db

import (
	"errors"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Product struct {
	Product_name        string `json: "product_name"`
	Product_description string `json: "product_description"`
	Product_price       int16  `json: "product_price"`
	Id_warehouse        uint8  `json:"id_warehouse"`
}

func GetAllProducts() ([]*Product, error) {

	rows, err := db.Query(`select name_product, description_product, price_product
						   from products;`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	products := make([]*Product, 0)

	for rows.Next() {
		product := new(Product)
		err := rows.Scan(
			&product.Product_name,
			&product.Product_description,
			&product.Product_price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}
func AddProduct(pr *Product) error {

	// Check current product.

	rows, err := db.Query(`select name_product, id_warehouse from Products where name_product = ? AND id_warehouse = ? `, &pr.Product_name, &pr.Id_warehouse)

	if err != nil {
		log.Fatal(err)
	}

	if !rows.Next() {
		_, err := db.Exec(`insert into products(id_warehouse,name_product,price_product,description_product) VALUES(?,?,?,?)`, &pr.Id_warehouse, &pr.Product_name, &pr.Product_price, &pr.Product_description)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	}
	return errors.New("This product already exists.")
}
