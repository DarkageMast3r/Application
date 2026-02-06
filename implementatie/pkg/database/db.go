package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"ZorgTechImplementatie/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	Offset(offset int) *gorm.DB
	Limit(limit int) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	Create(value interface{}) Database
	Where(query interface{}, args ...interface{}) Database
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Model(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) Database
	Updates(value interface{}) *gorm.DB
	Order(value interface{}) *gorm.DB
	Save(value interface{}) Database
	Error() error
}

type GormDatabase struct {
	*gorm.DB
}

func (db *GormDatabase) Create(value interface{}) Database {
	return &GormDatabase{db.DB.Create(value)}
}
func (db *GormDatabase) Where(query interface{}, args ...interface{}) Database {
	return &GormDatabase{db.DB.Where(query, args...)}
}

func (db *GormDatabase) First(dest interface{}, conds ...interface{}) Database {
	return &GormDatabase{db.DB.First(dest, conds...)}
}

func (db *GormDatabase) Error() error {
	return db.DB.Error
}

func (db *GormDatabase) Save(value interface{}) Database {
	return &GormDatabase{db.DB.Save(value)}
}

func NewDatabase() *gorm.DB {
	var database *gorm.DB
	var err error

	db_hostname := os.Getenv("POSTGRES_HOST")
	db_name := os.Getenv("POSTGRES_DB")
	db_user := os.Getenv("POSTGRES_USER")
	db_pass := os.Getenv("POSTGRES_PASSWORD")
	db_port := os.Getenv("POSTGRES_PORT")

	dbURl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db_user, db_pass, db_hostname, db_port, db_name)

	for i := 1; i <= 3; i++ {
		database, err = gorm.Open(postgres.Open(dbURl), &gorm.Config{})
		if err == nil {
			break
		} else {
			log.Printf("Attempt %d: Failed to initialize database. Retrying...", i)
			time.Sleep(3 * time.Second)
		}
	}

	// AutoMigrate voor alle modellen
	database.AutoMigrate(
		//&models.User{},
		&models.ImplementatieDossier{},
		&models.ZorgTechProduct{},
	)

	return database
}
