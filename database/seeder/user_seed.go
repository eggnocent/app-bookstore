package seeder

import (
	"app-bookstore/lib"
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func SeedUsers(db *sqlx.DB) {
	pass, err := lib.HashPassword("admin123")
	if err != nil {
		log.Fatal("Error hashing this password", err)
	}

	query := `
    INSERT INTO users (id, username, password, role, created_at, created_by, updated_at, updated_by)
    VALUES ($1, $2, $3, $4, NOW(), $5, NOW(), $5)
    ON CONFLICT (username) DO NOTHING;
    `

	adminID := uuid.New().String()

	_, err = db.ExecContext(context.Background(), query, adminID, "admin", pass, "admin", adminID)
	if err != nil {
		log.Fatalf("Error seeding users: %v", err)
	} else {
		log.Println("Users seeded successfully")
	}
}
