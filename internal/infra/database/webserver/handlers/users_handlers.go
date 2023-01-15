package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vitorconti/go-user-products-api/internal/dto"
	"github.com/vitorconti/go-user-products-api/internal/entity"
	"github.com/vitorconti/go-user-products-api/internal/infra/database"
)

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(userDB database.UserInterface) *UserHandler {
	return &UserHandler{UserDB: userDB}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := entity.NewUser(user.Name, user.Name, user.Password)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
	}
	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}
