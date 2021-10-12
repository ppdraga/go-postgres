package fixtures

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type R struct {
	DB *gorm.DB
	//conn   *sql.DB
}

func InitTestDB() (*R, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=cdn_test port=5432 TimeZone=Europe/Moscow"
	fmt.Println("dsn: ", dsn)
	dbcon, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err == nil {
		fmt.Println("Connected to DB!!!")
	} else {
		return nil, err
	}
	dbcon.Exec(`
		DROP TABLE IF EXISTS serverfile;
		DROP TABLE IF EXISTS userfile;
		DROP TABLE IF EXISTS "user";
		DROP TABLE IF EXISTS server;
		DROP TABLE IF EXISTS file;
		DROP TABLE IF EXISTS area;
		
		CREATE TABLE IF NOT EXISTS area (
			id serial PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			removed TIMESTAMP NULL,
			name VARCHAR(255) NOT NULL,
			description VARCHAR(255) NULL
			);
		
		CREATE TABLE IF NOT EXISTS file (
			id serial PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			removed TIMESTAMP NULL,
			name VARCHAR(255) NOT NULL,
			sha VARCHAR(255) NOT NULL,
			size BIGINT NOT NULL DEFAULT 0,
			description VARCHAR(255) NULL
			);
		
		CREATE TABLE IF NOT EXISTS server (
			id serial PRIMARY KEY,
			area_id BIGINT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			removed TIMESTAMP NULL,
			name VARCHAR(255) NOT NULL,
			hostname VARCHAR(255) NOT NULL,
			description VARCHAR(255) NULL,
			FOREIGN KEY (area_id) REFERENCES area(id) ON UPDATE CASCADE ON DELETE CASCADE
			);
		
		CREATE TABLE IF NOT EXISTS "user" (
			id serial PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			removed TIMESTAMP NULL,
			name VARCHAR(255) NOT NULL,
			balance BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0)
			);
		
		CREATE TABLE IF NOT EXISTS userfile (
			user_id BIGINT NOT NULL,
			file_id BIGINT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			removed TIMESTAMP NULL,
			FOREIGN KEY (user_id) REFERENCES "user"(id) ON UPDATE CASCADE ON DELETE CASCADE,
			FOREIGN KEY (file_id) REFERENCES file(id) ON UPDATE CASCADE ON DELETE CASCADE
			);
		
		CREATE TABLE IF NOT EXISTS serverfile (
			server_id BIGINT NOT NULL,
			file_id BIGINT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			removed TIMESTAMP NULL,
			FOREIGN KEY (server_id) REFERENCES server(id) ON UPDATE CASCADE ON DELETE CASCADE,
			FOREIGN KEY (file_id) REFERENCES file(id) ON UPDATE CASCADE ON DELETE CASCADE
			);
	`)
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
