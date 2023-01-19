package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	configs "github.com/vitorconti/go-user-products-api/config"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/vitorconti/go-user-products-api/docs"
	"github.com/vitorconti/go-user-products-api/internal/entity"
	"github.com/vitorconti/go-user-products-api/internal/infra/database"
	"github.com/vitorconti/go-user-products-api/internal/infra/database/webserver/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title           Go Expert API Example
// @version         1.0
// @description     Product API with auhtentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   Jose Vitor

// @host      localhost:8000
// @BasePath  /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	configs, err := configs.LoadConfig(".")
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

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("jwtExpiresIn", configs.JWTExpiresIn))
	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Get("/", productHandler.GetAllProducts)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
		r.Post("/users", userHandler.CreateUser)
	})
	r.Post("/users/generate_token", userHandler.GetJWT)
	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))
	http.ListenAndServe(":8000", r)
}
