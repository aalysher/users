package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"users/internal/models"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	CreateUser(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	UpdateUserByID(id string, updates models.UserUpdate) (*models.User, error)
}

type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}

func (s *service) CreateUser(user *models.User) error {
	id := uuid.New()
	query := `
        INSERT INTO users (id, first_name, last_name, email, age)
        VALUES ($1, $2, $3, $4, $5)
    `
	log.Printf("Executing query: %s with values: %s, %s, %s, %s, %d", query, id, user.FirstName, user.LastName, user.Email, user.Age)
	_, err := s.db.Exec(query, id, user.FirstName, user.LastName, user.Email, user.Age)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}
	//user.ID = id.String() // Устанавливаем ID в объекте user
	return nil
}

func (s *service) GetUserByID(id string) (*models.User, error) {
	var user models.User
	query := `SELECT id, first_name, last_name, email, age FROM users WHERE id = $1`
	err := s.db.QueryRow(query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Age)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *service) UpdateUserByID(id string, updates models.UserUpdate) (*models.User, error) {
	query := "UPDATE users SET "
	params := []interface{}{}
	paramId := 1

	if updates.FirstName != nil {
		query += fmt.Sprintf("first_name = $%d, ", paramId)
		params = append(params, *updates.FirstName)
		paramId++
	}
	if updates.LastName != nil {
		query += fmt.Sprintf("last_name = $%d, ", paramId)
		params = append(params, *updates.LastName)
		paramId++
	}
	if updates.Age != nil {
		query += fmt.Sprintf("age = $%d, ", paramId)
		params = append(params, *updates.Age)
		paramId++
	}
	if updates.Email != nil {
		query += fmt.Sprintf("email = $%d, ", paramId)
		params = append(params, *updates.Email)
		paramId++
	}

	// Remove the last comma and add the WHERE clause
	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d RETURNING id, first_name, last_name, email, age", paramId)
	params = append(params, id)

	var user models.User
	err := s.db.QueryRow(query, params...).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Age)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
