package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	url := os.Getenv("VCT_POSTGRES_URL")
	if url == "" {
		fmt.Println("VCT_POSTGRES_URL must be set")
		os.Exit(1)
	}

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Delete both the normal and the ".down" entry from schema_migrations
	tag, err := conn.Exec(context.Background(), "DELETE FROM schema_migrations WHERE version LIKE '0088%'")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Deleted %d rows from schema_migrations\n", tag.RowsAffected())
}
