package seeder

import (
	"app-bookstore/lib"
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// SeedSuperAdmin inserts only one Super Admin user.
func SeedSuperAdmin(db *sqlx.DB) {
	// Hash the password for Super Admin
	hashedPassword, err := lib.HashPassword("SuperAdmin123!")
	if err != nil {
		log.Fatalf("Error hashing Super Admin password: %v", err)
	}

	// Generate Super Admin ID
	superAdminID := uuid.New()

	// Insert Super Admin user
	query := `
		INSERT INTO users (id, username, password, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, NOW(), $4, NOW(), $4)
		ON CONFLICT (username) DO NOTHING;
	`

	_, err = db.ExecContext(context.Background(), query, superAdminID, "superadmin", hashedPassword, superAdminID)
	if err != nil {
		log.Fatalf("Error inserting Super Admin: %v", err)
	} else {
		log.Println("Super Admin user inserted successfully.")
	}
}
