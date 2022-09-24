package auth

import (
	"assignment2/constant"
	"assignment2/entity"
	"assignment2/user"
	"assignment2/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type authHandler struct {
	r  *mux.Router
	as user.AuthService
}

func ProvideAuthHandler(r *mux.Router, as user.AuthService) *authHandler {
	return &authHandler{r: r, as: as}
}

func (a *authHandler) InitHandler() {
	route := a.r.PathPrefix(constant.AUTH_USER_API_PATH).Subrouter()

	route.HandleFunc("/register", a.newUser()).Methods(http.MethodPost)
	route.HandleFunc("/login", a.loginUser()).Methods(http.MethodPost)
}

func (a *authHandler) newUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := entity.UserRegistrationRequest{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Printf("[newUser] failed to parse JSON data, err => %v", err)
		}
		res, err := a.as.NewUser(r.Context(), &data)
		if err != nil {
			log.Printf("[newUser] failed to store new user, err => %v", err)
		}
		utils.NewBaseResponse(http.StatusCreated, "SUCCESS", nil, res).SendResponse(&w)
	}
}

func (a *authHandler) loginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := entity.UserLoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			log.Printf("[loginUser] failed to parse the JSON data, err => %v", err)
			return
		}

		res, err := a.as.LoginUser(r.Context(), &data)
		if err != nil {
			log.Printf("[loginUser] failed to login, err => %v", err)
			return
		}
		utils.NewBaseResponse(http.StatusOK, "SUCCESS", nil, res).SendResponse(&w)
	}
}
