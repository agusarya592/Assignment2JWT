package main

import (
	"assignment2/config"
	"assignment2/database"
	"assignment2/route"
	"assignment2/server"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {

	//localhost:8080/api/orders
	config.Init()
	cfg := config.GetConfig()
	db := database.Init()
	root := mux.NewRouter()
	route.Init(root, db)
	s := server.ProvideServer(cfg.ServerAddress, root)
	s.ListenAndServe()
}
