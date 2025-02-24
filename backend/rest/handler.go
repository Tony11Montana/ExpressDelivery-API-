package rest

import (
	or "backend/models_db"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func Hi(w http.ResponseWriter, r *http.Request) {

	//w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "text/plain")

	orders, err := or.GetAllOrder()

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	ors, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(ors)
	//fmt.Fprintf(w, "Hi, server!")

}
