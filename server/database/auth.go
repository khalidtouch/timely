package database

import (
	"encoding/json"
	"log"
)

func Login(req []byte) User {
	var user User 
	db := prepareDb(dbname)
	defer db.Close() 

	json.Unmarshal(req, &user)
	temp := getUserByUserNameAndPassword(user.UserName, user.Password)
	if temp.UserName != user.UserName || temp.Password != user.Password {
		log.Print("Please provide the correct credentials")
		user = User{}
	} else {
		addTokenToUser(&user)
		saveTokenToDb(&user)
	}
	return user 
}


func Logout(username string) (int64, error) {
	db := prepareDb(dbname)
	defer db.Close() 
	user := getUserByUserName(username)

	stmt, err := db.Prepare("delete from Token where UserId = ?")
	if err != nil {
		log.Print("An error occurred while deleting token")
		return 0, err 
	}

	defer stmt.Close() 

	res, err := stmt.Exec(user.Id)
	if err != nil {
		log.Print("An error occurred while deleting all tokens belong to user ", err)
		return 0, err 
	}

	return res.RowsAffected()
}

