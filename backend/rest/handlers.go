package rest

import (
	or "backend/models_db"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

var jwtKey = []byte("Elagin_diplom")

func ParseJWTToken(tokenString string, signingKey []byte) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &or.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing mwthod : %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(*or.Claims); ok && token.Valid {
		return claims.Login_user, claims.Role, nil
	}
	return "", "", err
}

func AllOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
		return
	}

	tokenString := parts[1]

	login, role, err := ParseJWTToken(tokenString, jwtKey)

	if err != nil {
		http.Error(w, "Invalid authorization (JWT token end or not use)", http.StatusUnauthorized)
		return
	}

	orders, err := or.GetAllOrder(&login, &role)

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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &or.Claims{
		Login_user: user.Login_user,
		Role:       role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
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

	/*check, err, role := or.CheckUser(&usr)
	if err != nil || !check {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	*/

	// Генерация JWT-токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &or.Claims{
		Login_user: usr.Login_user,
		Role:       usr.Role_user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
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
