package Database

import (
	"216/internal/orchestrator/Entities"
	"errors"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Instance *gorm.DB
var err error

func Connect(connectionString string) {
	log.Println("Connected to Database...")
	Instance, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		panic("Cannot connect to DB")
	}
	log.Println("Success")
}

func Migrate() {
	Instance.AutoMigrate(
		&Entities.ArithmeticOperation{},
		&Entities.ArithmeticExpressions{},
		&Entities.ComputingResource{},
		&Entities.User{},
	)
	log.Println("Database Migration Completed...")
}

func Seeder() {
	operations := [4]string{"+", "-", "*", "/"}
	Instance.Exec("DELETE FROM arithmetic_operations")
	if Instance.Migrator().HasTable(&Entities.ArithmeticOperation{}) {
		if err := Instance.First(&Entities.ArithmeticOperation{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Seeding...")
			for _, v := range operations {
				user := Entities.ArithmeticOperation{Value: v, LeadTime: 10}
				Instance.Create(&user)
			}
			log.Println("Success")
		}
	}
}
