package main

import (
	db "backend/models_db"
	"backend/rest"
	"net/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db.InitDB("root:admin@tcp(127.0.0.1:3306)/ElaginDiplom")
	
	router := mux.NewRouter()
	router.HandleFunc("/orders", rest.Hi).Methods("GET")

	http.ListenAndServe(":80",
	handlers.CORS(
		handlers.AllowedOrigins([]string {"*"}),
		handlers.AllowedMethods([]string {"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string {"X-Requested-With", "Content-Type", "Authorization"}),

	)(router))
}

/* db, err := sql.Open("mysql", "root:admin@tcp(database:3306)/myTest")

if err != nil {
	panic(err)
}
defer db.Close()

database = db

name := "phone"

_, err = database.Exec("insert into myTest.Products (name) values (?);",
	name)

if err != nil {
	fmt.Println("Dont try insert in table")
	panic(err)
}
*/
