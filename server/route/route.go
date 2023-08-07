package route

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/dgrijalva/jwt-go"
	data "github.com/khalidtouch/timely/database"
)



var counter int 
var mutex = &sync.Mutex{}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []data.User = data.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	req, _ := io.ReadAll(r.Body)
	user := data.RegisterUser(req)

	if (data.User{}) == user {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode("This username already exists")
	} else {
		r.Header.Set("Token", user.Token)
		json.NewEncoder(w).Encode(r.Header)
	}
}


func Login(w http.ResponseWriter, r *http.Request) {
	req, _ := io.ReadAll(r.Body)
	user := data.Login(req)

	if (data.User{}) == user {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Incorrect credentials")
	} else {
		r.Header.Set("Token", user.Token)
		json.NewEncoder(w).Encode(r.Header)
		w.WriteHeader(http.StatusCreated)
		fmt.Println("author ---> ", r.Header.Get("Authorization"))
	}
}

func GetOwnAccount(w http.ResponseWriter, r *http.Request) {
	props := r.Context().Value("props")
	if mapper, ok := props.(jwt.MapClaims); ok {
		username := mapper["user_name"] 
		if v, e := username.(string); e {
			if v == "" {
				json.NewEncoder(w).Encode("Please login again")
			}

			user := data.GetMe(v)
			user.Token = data.GetLastLoginToken(v)
			if user.UserName == "" || user.Token == "" || len(data.CheckToken(user.Token)) == 0 {
				r.Header.Set("Authorization", "")
				fmt.Println("del uth --> ", r.Header.Get("Authorization"))
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode("Please authenticate")
			} else {
				json.NewEncoder(w).Encode(user)
			}
		}
	}
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	props := r.Context().Value("props")
	req, _ := io.ReadAll(r.Body)

	if mapper, ok := props.(jwt.MapClaims); ok {
		username := mapper["user_name"]
		if v, e := username.(string); e {
			user := data.GetMe(v)
			if user.Token == "" || user.UserName == "" {
				w.WriteHeader(http.StatusUnauthorized)
				r.Header.Set("Authorization", "")
				json.NewEncoder(w).Encode("Please authenticate")
				return 
			}

			data.CreateTask(v, req)
			json.NewEncoder(w).Encode("Your task has been saved successfully")
		}
	} else {
		json.NewEncoder(w).Encode("Unable to create Task")
	}
}


func GetTasks(w http.ResponseWriter, r *http.Request) {
	props := r.Context().Value("props")
	if mapper, ok := props.(jwt.MapClaims); ok {
		username := mapper["user_name"]
		if v, e := username.(string); e {
			tasks := data.GetTasksUnderUser(v)
			json.NewEncoder(w).Encode(tasks)
		}
	}
}

func DeleteOwnAccount(w http.ResponseWriter, r *http.Request) {
	props := r.Context().Value("props")
	if mapper, ok := props.(jwt.MapClaims); ok {
		username := mapper["user_name"]
		if v, e := username.(string); e {
			if id := data.DeleteMe(v); id != 0 {
				json.NewEncoder(w).Encode("Your account has been deleted successfully")
				return 
			}
			w.WriteHeader(http.StatusNotFound)
			log.Print("An error occured while deleting task")
		}
	}
}


func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	props := r.Context().Value("props")
	req, _ := io.ReadAll(r.Body)
	if mapper, ok := props.(jwt.MapClaims); ok {
		username := mapper["user_name"]
		if v, e := username.(string); e {
			w.WriteHeader(http.StatusConflict)
			if user, err := data.UpdateUser(req, v); err != nil {
				fmt.Fprintf(w, "Failed to update user!")
			} else {
				r.Header.Set("Token", user.Token)	
				json.NewEncoder(w).Encode(r.Header)	
			}
		
		}
	}
}


func Logout(w http.ResponseWriter, r *http.Request) {
	props := r.Context().Value("props")
	if mapper, ok := props.(jwt.MapClaims); ok {
		username := mapper["user_name"]
		if v, e := username.(string); e {
			_, err := data.Logout(v)
			if err != nil {
				log.Print("An error occured while logging out ", err)
			}
			user := data.GetMe(v)
			authR := strings.Split(r.Header.Get("Authorization"), "Bearer ")
			user.Token = authR[1]
			r.Header.Set("Token", user.Token)
			json.NewEncoder(w).Encode(r.Header)
		}
	}
}