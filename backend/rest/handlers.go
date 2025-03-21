package rest

import (
	or "backend/models_db"
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

var jwtKey = []byte("Elagin_diplom")

func AllOrder(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Content-Type", "text/plain")

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

func AllCouriers(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Content-Type", "text/plain")

	couriers, err := or.GetAllCouriers()

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
	couriersJSON, err := json.Marshal(couriers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(couriersJSON)
	//fmt.Fprintf(w, "Hi, server!")

}

func AddCourier(w http.ResponseWriter, r *http.Request) {
	var courier or.Courier
	err := json.NewDecoder(r.Body).Decode(&courier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = or.AddCourier(&courier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Courier added successfully"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user or.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка учетных данных
	check, err, role := or.CheckUser(&user)
	if err != nil || !check {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Генерация JWT-токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":    user.Login_user,
		"password": user.Password_user,
		"role":     role,
		//	"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var usr or.User

	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = or.AddUser(&usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	check, err, role := or.CheckUser(&usr)
	if err != nil || !check {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Генерация JWT-токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":    usr.Login_user,
		"password": usr.Password_user,
		"role":     role,
		//	"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
func AllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	products, err := or.GetAllProducts()

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	prods, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(prods)
}
func AddProduct(w http.ResponseWriter, r *http.Request) {
	var product or.Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = or.AddProduct(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

}
