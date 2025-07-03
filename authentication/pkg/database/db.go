package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"authentication/pkg/models" // Zorg dat dit pad klopt

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database defines the consistent interface for database operations
// database/db.go

// Database defines the consistent interface for database operations
type Database interface {
	Offset(offset int) Database
	Limit(limit int) Database
	Find(dest interface{}, conds ...interface{}) Database
	Create(value interface{}) Database
	Where(query interface{}, args ...interface{}) Database
	Delete(value interface{}, conds ...interface{}) Database
	Model(value interface{}) Database
	First(dest interface{}, conds ...interface{}) Database
	Updates(value interface{}) Database
	Update(column string, value interface{}) Database // <--- ADD THIS LINE
	Order(value interface{}) Database
	Save(value interface{}) Database
	Preload(query string, args ...interface{}) Database
	Error() error
	WithContext(ctx context.Context) Database
	Raw(sql string, values ...interface{}) Database
	Scan(dest interface{}) Database
	Exec(sql string, values ...interface{}) Database
}

// GormDatabase wraps *gorm.DB to implement the Database interface
type GormDatabase struct {
	*gorm.DB
}

// Implement methods to return `Database` interface consistently
func (db *GormDatabase) Offset(offset int) Database { return &GormDatabase{db.DB.Offset(offset)} }
func (db *GormDatabase) Limit(limit int) Database   { return &GormDatabase{db.DB.Limit(limit)} }
func (db *GormDatabase) Find(dest interface{}, conds ...interface{}) Database {
	return &GormDatabase{db.DB.Find(dest, conds...)}
}
func (db *GormDatabase) Create(value interface{}) Database { return &GormDatabase{db.DB.Create(value)} }
func (db *GormDatabase) Where(query interface{}, args ...interface{}) Database {
	return &GormDatabase{db.DB.Where(query, args...)}
}
func (db *GormDatabase) Delete(value interface{}, conds ...interface{}) Database {
	return &GormDatabase{db.DB.Delete(value, conds...)}
}
func (db *GormDatabase) Model(value interface{}) Database { return &GormDatabase{db.DB.Model(value)} }
func (db *GormDatabase) First(dest interface{}, conds ...interface{}) Database {
	return &GormDatabase{db.DB.First(dest, conds...)}
}

func (db *GormDatabase) Updates(value interface{}) Database {
	return &GormDatabase{db.DB.Updates(value)}
}
func (db *GormDatabase) Update(column string, value interface{}) Database { // <--- ADD THIS METHOD
	return &GormDatabase{db.DB.Update(column, value)}
}
func (db *GormDatabase) Order(value interface{}) Database { return &GormDatabase{db.DB.Order(value)} }
func (db *GormDatabase) Save(value interface{}) Database  { return &GormDatabase{db.DB.Save(value)} }

func (db *GormDatabase) Preload(query string, args ...interface{}) Database {
	return &GormDatabase{db.DB.Preload(query, args...)}
}

//	func (db *GormDatabase) Association(column string) *gorm.Association {
//		return db.DB.Association(column)
//	}
func (db *GormDatabase) Error() error { return db.DB.Error }
func (db *GormDatabase) WithContext(ctx context.Context) Database {
	return &GormDatabase{db.DB.WithContext(ctx)}
}

// Oplossing voor de Exec, Raw, en Scan methoden:
func (db *GormDatabase) Raw(sql string, values ...interface{}) Database {
	return &GormDatabase{db.DB.Raw(sql, values...)}
}

func (db *GormDatabase) Scan(dest interface{}) Database {
	return &GormDatabase{db.DB.Scan(dest)}
}

func (db *GormDatabase) Exec(sql string, values ...interface{}) Database {
	return &GormDatabase{db.DB.Exec(sql, values...)}
}

// NewGormDatabase initializes a new GORM database connection and runs migrations.
// It returns the custom Database interface.
func NewGormDatabase() (Database, error) {
	var database *gorm.DB
	var err error

	dbHostname := os.Getenv("POSTGRES_HOST")
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbPort := os.Getenv("POSTGRES_PORT")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHostname, dbPort, dbName) // Typo: dbURl naar dbURL

	for i := 1; i <= 3; i++ {
		database, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{}) // Gebruik dbURL
		if err == nil {
			log.Println("Database connection successful.")
			break
		} else {
			log.Printf("Attempt %d: Failed to initialize database. Retrying... Error: %v", i, err)
			time.Sleep(3 * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after multiple retries: %w", err)
	}

	// AutoMigrate all models
	err = database.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.APIEndpoint{},
		&models.RefreshToken{}, // Changed from RefreshToken to AuthToken based on prior discussion
		&models.PasswordResetToken{},
		&models.TokenBlacklist{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate database: %w", err)
	}

	return &GormDatabase{database}, nil // Return the custom interface
}
