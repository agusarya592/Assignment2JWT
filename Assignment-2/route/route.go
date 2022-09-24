package route

import (
	"assignment2/auth"
	"assignment2/controller"
	"assignment2/middleware"
	"assignment2/repository"
	"assignment2/service"
	"assignment2/user"
	"database/sql"

	"github.com/gorilla/mux"
)

//localhost:8080/api/orders

func Init(router *mux.Router, db *sql.DB) {
	webRouter := router.NewRoute().PathPrefix("/api").Subrouter()
	protectedRoute := router.NewRoute().PathPrefix("/api").Subrouter()

	protectedRoute.Use(middleware.AuthMiddleware())

	orderRepository := repository.ProvideRepository(db)
	orderService := service.ProvideService(orderRepository)
	orderHandler := controller.ProvideController(webRouter, protectedRoute, orderService)
	orderHandler.InitController()

	userRepository := user.ProvideAuthRepository(db)
	userService := user.ProvideAuthService(userRepository)
	userHandler := auth.ProvideAuthHandler(webRouter, userService)
	userHandler.InitHandler()
}
