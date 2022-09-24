package database

import (
	"assignment2/config"
	"database/sql" //sql golang *sql.DB
	"fmt"
	"log"
)

func Init() *sql.DB {
	cfg := config.GetConfig()
	log.Printf("TEST")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Address, cfg.Database.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("ERROR Connection to db, err: %v", err)
	}
	log.Printf("SUCCESS")
	return db
}
