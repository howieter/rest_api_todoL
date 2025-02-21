package todoAPI

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func CreateDB() error {
	dbURL := getPostgresURL()
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("error while connecting to postgres DB: %v", err)
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "CREATE DATABASE todoDB")
	if err != nil {
		return fmt.Errorf("error while creating todoDB: %v", err)
	}
	return nil
}

func ConnectDB() (*pgx.Conn, error) {
	dbURL := getPostgresURL()
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to DB: %v", err)
	}
	return conn, err
}

func getPostgresURL() string {
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}
	return fmt.Sprintf("postgres://%s@localhost:5432/tododb?sslmode=disable", user)
}

// todo: make migrates with pgx
func MakeQuery(conn *pgx.Conn, filePath string) error {
	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error while reading file %s: %w", filePath, err)
	}
	sqlContent := string(sqlBytes)

	_, err = conn.Exec(context.Background(), sqlContent)
	if err != nil {
		return fmt.Errorf("error while making sql query %s: %w", filePath, err)
	}

	return nil
}
