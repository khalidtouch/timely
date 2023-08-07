package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	environment "github.com/joho/godotenv"
)

type (
	User struct {
		Id		int `json:"UserId"`
		UserName string `json:"UserName"`
		Email	string 	`json:"Email"`
		Password string `json:"Password"`
		Token 	 string `json:"Token"`
	}

	Token struct {
		UserId	int	 `json:"UserId"`
		Value string `json:"Value"`
	}

	Task struct {
		UserId	int `json:"UserId"`
		Completed bool `json:"Completed"`
		Name 	string `json:"Name"`
	}
)


var db *sql.DB 

func init() {
	err := environment.Load()
	if err != nil {
		log.Fatal("Error loading .env file; ", err)
	}


	db = prepareDb()
	defer db.Close()
}


func prepareDb() *sql.DB {
	connection := os.Getenv("DB_CONNECTION")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	hostname := os.Getenv("DB_HOSTNAME")
	dbname := os.Getenv("DB_NAME")

	db, err := sql.Open(connection, dsn(username, password, hostname))
	if err != nil {
		log.Printf("error %s during the open db\n", err)
	}

	connectToDb(dbname, db)
	return db 
}


func dsn(username, password, hostname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/", username, password, hostname)
}


func connectToDb(dbname string, db *sql.DB) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS " + dbname)
	if err != nil {
		log.Printf("error %s when creating db\n", err)
		return 
	}

	_, err = db.ExecContext(ctx, "use " + dbname)
	if err != nil {
		log.Printf("error using db, %s", err)
		return 
	}

	_, err = db.ExecContext(ctx, "create table if not exists User (UserId int not null, UserName varchar(50) default null, Email varchar(50) default null, Password varchar(50) default null, Token text default null, primary key(UserId))")
	if err != nil {
		log.Printf("error creating table User, %s", err)
		return 
	}

	_, err = db.ExecContext(ctx, "create table if not exists Token (UserId int default null, Value text default null)")
	if err != nil {
		log.Printf("error creating table Token, %s", err)
		return 
	}


	_, err = db.ExecContext(ctx, "create table if not exists Task (UserId int default null, Completed boolean default false, Name text default null)")
	if err != nil {
		log.Printf("error creating table Task, %s", err)
		return 
	}

	no, err := res.RowsAffected() 
	if err != nil {
		log.Printf("error %s when fetching db\n", err)
	}
	log.Printf("rows affected %d\n", no)	
}

