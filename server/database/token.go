package database

import (
	"github.com/khalidtouch/timely/middleware"
	"log"
	"errors"
)


func addTokenToUser(user *User) {
	userId := getUserIdByName(user.UserName)

	token, err := middleware.GenerateToken(uint64(userId), user.UserName)
	if err != nil {
		log.Fatal("An error occured during the produce token ", err)
	}

	user.Token = token 
	updateUserOnDb(user, int64(userId))
}


func saveTokenToDb(user *User) (int64, error) {
	userId := getUserIdByName(user.UserName)
	var token = Token{UserId: userId, Value: user.Token}

	if token.Value == "" {
		log.Print("Invalid Token")
		return 0, errors.New("Invalid token error")
	}

	db := prepareDb()
	defer db.Close() 

	stmt, err := db.Prepare("insert into Token values (?,?)")
	if err != nil {
		log.Fatal("An error occured while inserting token ", err)
		return 0, err 
	}

	res, err := stmt.Exec(token.UserId, token.Value)
	if err != nil {
		log.Fatal("An error occured while adding token to db ", err)
		return 0, err 
	}
	return res.RowsAffected()
}


func GetLastLoginToken(username string) string {
	db := prepareDb()
	defer db.Close() 

	user := getUserByUserName(username)
	res, err := db.Query("select Value from Token where UserId = ?", user.Id)
	if err != nil {
		log.Print("An error occured while retrieving the last token", err)
		return ""
	}

	var tokens []string 
	for res.Next() {
		var temp Token 
		err := res.Scan(&temp.Value)
		if err != nil {
			log.Print("An error occured while scanning db to get token")
			break 
		}

		tokens = append(tokens, temp.Value)
	}

	if len(tokens) == 0 {
		return ""
	} else {
		return tokens[len(tokens) - 1]
	}
}


func CheckToken(token string) []string {
	var tokens []string 
	db := prepareDb()
	defer db.Close()

	res, err := db.Query("select * from Token where Value = ?", token)
	if err != nil {
		log.Print("An error occured while fetching token from db", err)
		return make([]string, 0)
	} else {
		for res.Next() {
			var temp Token 
			err = res.Scan(&temp.UserId, &temp.Value)
			if err != nil {
				log.Print("An error while scanning db to get token")
				return make([]string, 0)
			} else {
				tokens = append(tokens, temp.Value)
			}
		}
	}
	return tokens 
}


func deleteAllTokens(user *User) (int64, error) {
	db := prepareDb()
	defer db.Close() 

	stmt, err := db.Prepare("delete from Token where UserId = ?")
	if err != nil {
		log.Print("An error occured while deleting tokens", err)
		return 0, err 
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.Id)
	if err != nil {
		log.Print("An error occured while deleting db")
		return 0, err 
	}
	return res.RowsAffected()
}