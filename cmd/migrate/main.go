package main

import (
	"log"
	"kerjadekat/backend/config"
	"kerjadekat/backend/internal/domain"
	"kerjadekat/backend/pkg/database"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Config load failed: %v", err)
	}

	db, err := database.OpenPostgres(cfg)
	if err != nil {
		log.Fatalf("DB Open Error: %v", err)
	}

	log.Println("Running AutoMigrate for User schema...")
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("AutoMigrate Error: %v", err)
	}
	
	log.Println("Database migration completed successfully!")
}
