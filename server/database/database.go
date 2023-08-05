package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

type (
	User struct {
		Id		int `json:"UserId"`
		UserName string `json:"UserName"`
		Email	string 	`json:"Email"`
		Password string `json:"Password"`
		CreatedAt string `json:"CreatedAt"`
	}

	Token struct {
		UserId	int	 `json:"UserId"`
		Value string `json:"Value"`
	}

	Task struct {
		UserId	int `json:"UserId"`
		Completed bool `json:"Completed"`
		Name 	string `json:"Name"`
		CreatedAt string `json:"CreateAt"`
	}
)

const (
	username = "root"
	password = "IDeyTellYou555!"
	hostname = "127.0.0.1:3306"
	dbname = "timely_db"
)


func init() {

}


func dsn(databaseName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, databaseName)
}


func connectToDb(db *sql.DB) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS " + dbname)
	
	if err != nil {
		log.Printf("error %s when creating db\n", err)
		return 
	}

	no, err := res.RowsAffected() 
	if err != nil {
		log.Printf("error %s when fetching db\n", err)
	}
	log.Printf("rows affected %d\n", no)	
}


func prepareDb(dbName string) *sql.DB {
	db, err := sql.Open("mysql", dsn(dbName))
	if err != nil {
		log.Printf("error %s during the open db\n", err)
	}
	connectToDb(db)
	return db 
}


func GetAllUsers() []User 

