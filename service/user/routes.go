package user

import (
	"fmt"
	"github.com/ViniciusDSLima/golang01/config"
	"github.com/ViniciusDSLima/golang01/service/auth"
	"github.com/ViniciusDSLima/golang01/types"
	"github.com/ViniciusDSLima/golang01/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/user/{userId}", auth.WithJWTAuth(h.handlerGetUser, h.store)).Methods(http.MethodGet)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user types.LoginUserPayload

	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid payload: %v", errors))
		return
	}

	u, err := h.store.GetUsersByEmail(user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("User not found, invalid email"))
		return
	}

	if !auth.ComparePassword(u.Password, []byte(user.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid email or password"))
		return
	}

	secret := []byte(config.Env.JWTSecret)

	token, err := auth.CreateJWT(secret, u.Id)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user types.RegisterUserPayload

	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid payload: %v", errors))
		return
	}

	_, err := h.store.GetUsersByEmail(user.Email)

	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("User already exists", user.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	responsePayload := types.UserResponsePayload{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	utils.WriteJSON(w, http.StatusCreated, responsePayload)

}

func (h *Handler) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userId, ok := vars["userId"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid user id"))
		return
	}

	user, err := h.store.GetUserById(userId)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)

}
