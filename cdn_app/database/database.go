package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type R struct {
	DB *gorm.DB
	//conn   *sql.DB
}

func InitDB() (*R, error) {
	host := os.Getenv("PG_HOST")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWD")
	dbname := os.Getenv("PG_DB")
	port := os.Getenv("PG_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Europe/Moscow",
		host, user, password, dbname, port)
	fmt.Println("dsn: ", dsn)
	dbcon, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err == nil {
		fmt.Println("Connected to DB!!!")
	} else {
		return nil, err
	}
	return &R{DB: dbcon}, nil
}

func (r *R) Release() error {
	sqlDB, err := r.DB.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	if err != nil {
		return err
	}
	return nil
}
