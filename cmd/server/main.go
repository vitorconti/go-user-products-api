package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	configs "github.com/vitorconti/go-user-products-api/config"
	"github.com/vitorconti/go-user-products-api/internal/dto"
	"github.com/vitorconti/go-user-products-api/internal/entity"
	"github.com/vitorconti/go-user-products-api/internal/infra/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {

		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("teste.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})
	productDB := database.NewProduct(db)
	productHandler := NewProductHandler(productDB)
	fmt.Println("listen and serve")
	http.HandleFunc("/products", productHandler.CreateProduct)
	http.ListenAndServe(":8000", nil)

}

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
