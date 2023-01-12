package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	configs "github.com/vitorconti/go-user-products-api/config"

	"github.com/vitorconti/go-user-products-api/internal/entity"
	"github.com/vitorconti/go-user-products-api/internal/infra/database"
	"github.com/vitorconti/go-user-products-api/internal/infra/database/webserver/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {

		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("teste.sqlite"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})
	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/products", productHandler.CreateProduct)
	r.Get("/products/{id}", productHandler.GetProduct)
	http.ListenAndServe(":8000", r)
}
