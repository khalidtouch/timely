package database

import (
	"encoding/json"
	"log"
	"fmt"
	"database/sql"
)


func GetMe(username string) User {
	user := getUserByUserName(username)
	return user 
}


func UpdateUser(req []byte, username string) (User, error) {
	db := prepareDb(dbname)
	defer db.Close() 

	var newUser User 
	json.Unmarshal(req, &newUser)
	newUser.Id = getUserIdByName(username)
	_, err := updateUserOnDb(&newUser, int64(getUserIdByName(username)))
	if err != nil {
		log.Fatal("Failed to update user", err)
		return User{}, err 
	}
	addTokenToUser(&newUser)
	saveTokenToDb(&newUser)
	
	return newUser, nil 
}

func DeleteMe(username string) int64 {
	db := prepareDb(dbname)
	defer db.Close() 

	user := getUserByUserName(username)
	id, err := deleteUserById(db, int64(user.Id))
	if err != nil {
		log.Print("Failed to delete user from db", err)
	}


	deleteAllTokens(&user)
	deleteAllTasks(&user)
	log.Printf("Deleted row with ID %d\n", id)
	return id 
}


func deleteUserById(db *sql.DB, id int64) (int64, error) {
	stmt, err := db.Prepare("delete from User where UserId = ?")
	if err != nil {
		log.Print("An error occured while deleting from user")
		return 0, err 
	}
	defer stmt.Close() 

	res, err := stmt.Exec(id)
	if err != nil {
		log.Print("An error occured while deleting user")
		return 0, err 
	}
	return res.RowsAffected()
}


func GetAllUsers() []User {
	db := prepareDb(dbname)
	defer db.Close() 

	results, err := db.Query("select * from User")
	if err != nil {
		panic(err.Error())
	}

	var users []User 

	for results.Next() {
		var temp User 
		err = results.Scan(&temp.Id, &temp.UserName, &temp.Password, &temp.Token)
		if err != nil {
			panic(err.Error())
		}

		GetLastLoginToken(temp.UserName) //assumed redundancy
		users = append(users, temp)
	}

	return users 
}


func RegisterUser(req []byte) User {
	var user User 
	db := prepareDb(dbname)
	defer db.Close() 

	json.Unmarshal(req, &user)
	addTokenToUser(&user)

	id, err := saveUserToDb(db, &user)
	if err != nil {
		log.Println("Failed to save user into db", err)
		user = User{}
	}

	

	_, err = saveTokenToDb(&user)
	if err != nil {
		log.Println("An error occured while saving token to db ", err)
	}

	log.Printf("Inserted row with ID of %d\n", id)
	return user 
}



func getUserIdByName(username string) int {
	db := prepareDb(dbname)
	defer db.Close() 
	
	results, err := db.Query("select UserId from User where UserName = ?", username)
	if err != nil {
		log.Fatal("An error occured during the query db to get id by name ", err)
	}

	var id int 
	for results.Next() {
		err = results.Scan(&id)
		if err != nil {
			log.Fatal("An error occured during the scan db to get id by name", err)
		}
	}

	fmt.Println("got id ---> ", id)
	return id 
}


func getUserByUserName(username string) User {
	var user User 
	db := prepareDb(dbname)
	defer db.Close() 

	res, err := db.Query("select * from User where UserName = ?", username)
	if err != nil {
		log.Fatal("An error occured while fetching user")
	}

	for res.Next() {
		err := res.Scan(&user.Id, &user.UserName, &user.Email, &user.Password, &user.Token)
		if err != nil {
			log.Fatal("An error occured while scanning db to get user")
		}
	}
	return user 
}


func getUserByUserNameAndPassword(username, password string) User {
	var user User 
	db := prepareDb(dbname)
	defer db.Close()

	res, err := db.Query("select * from User where UserName = ? and Password = ?", username, password)
	if err != nil {
		log.Fatal("An error occured while fetching user by username and password")
	}

	for res.Next() {
		err := res.Scan(&user.Id, &user.UserName, &user.Email, &user.Password, &user.Token)
		if err != nil {
			log.Fatal("An error occured while reading db")
		}
	}
	return user 
}


func saveUserToDb(db *sql.DB, user *User) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO User VALUES (?,?,?,?,?)")

	if err != nil {
		return -1, err 
	}
	defer stmt.Close() 

	res, err := stmt.Exec(user.Id, user.UserName, user.Email, user.Password, user.Token)
	if err != nil {
		return -1, err 
	}
	return res.LastInsertId() 
}


func updateUserOnDb(user *User, userId int64) (int64, error) {
	db := prepareDb(dbname)
	defer db.Close() 
	stmt, err := db.Prepare("update User set UserName = ?, Email = ?, Password = ?, Token = ? where UserId = ?")
	if err != nil {
		log.Fatal("An error occured while updating user", err)
		return 0, err 
	}
	defer stmt.Close() 


	res, err := stmt.Exec(user.UserName, user.Email, user.Password, user.Token, userId)
	fmt.Println("new token ", user.Token)
	if err != nil {
		log.Fatal("An error occured while updating user", err)
		return 0, err 
	}
	return res.RowsAffected()
}