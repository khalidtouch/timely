package database

import (
	"encoding/json"
	"log"
)


func GetAllTasks() []Task {
	var tasks []Task 
	db := prepareDb(dbname)
	defer db.Close() 

	res, err := db.Query("select * from Task")
	if err != nil {
		log.Print("An error occured while getting task from db")
	}

	for res.Next() {
		var temp Task 
		err := res.Scan(&temp.UserId, &temp.Completed, &temp.Name)
		if err != nil {
			log.Print("An error occured while scanning db to get tasks")
		}
		tasks = append(tasks, temp)
	}
	return tasks 
}

func GetTasksUnderUser(username string) []string {
	db := prepareDb(dbname)
	defer db.Close() 

	user := getUserByUserName(username)
	res, err := db.Query("select * from Task where UserId = ?", user.Id)
	if err != nil {
		log.Fatal("An error occured while getting task from db")
	}

	var tasks []string 
	for res.Next() {
		var temp Task 
		err := res.Scan(&temp.UserId, &temp.Completed, &temp.Name)
		if err != nil {
			log.Fatal("An error occured while scanning db to get tasks")
		}
		tasks = append(tasks, temp.Name)
	}
	return tasks 
}


func CreateTask(username string, req []byte) User {
	user := getUserByUserName(username)
	var task Task 

	json.Unmarshal(req, &task)
	task.UserId = user.Id 

	_, err := saveTaskToDb(user, task)
	if err != nil {
		log.Fatal("An error occured while creating task for User", err)
	}
	return user 
}


func saveTaskToDb(user User, task Task) (int64, error) {
	db := prepareDb(dbname)
	defer db.Close() 

	stmt, err := db.Prepare("insert into Task VALUES(?,?,?)")
	if err != nil {
		log.Fatal("An error while saving task")
		return 0, err 
	}

	res, err := stmt.Exec(task.UserId, task.Completed, task.Name)
	if err != nil {
		log.Fatal("An error while saving task", err)
	}
	return res.RowsAffected()
}

func deleteAllTasks(user *User) (int64, error) {
	db := prepareDb(dbname)
	defer db.Close() 

	stmt, err := db.Prepare("delete from Task where UserId = ?")
	if err != nil {
		log.Print("An error occured while preparing to delete task")
		return 0, err 
	}
	defer stmt.Close() 

	res, err := stmt.Exec(user.Id)
	if err != nil {
		log.Print("An error occured while deleting task", err)
		return 0, err 
	}
	return res.RowsAffected()
}
