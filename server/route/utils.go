package route

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/khalidtouch/timely/middleware"
)

func IncrementCounter(w http.ResponseWriter, r *http.Request) {
	mutex.Lock() 
	counter++ 
	fmt.Fprintf(w, "Counter : %d", counter)
	mutex.Unlock()
}


func HandleRequest() {
	router := mux.NewRouter().StrictSlash(true)
	
	router.Handle("/", http.FileServer(http.Dir("./static")))
	router.HandleFunc("/counter", IncrementCounter)
	router.HandleFunc("/users", GetUsers)
	router.HandleFunc("/signUp", SignUp).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")
	router.Handle("/user", middleware.Authenticator(http.HandlerFunc(GetOwnAccount)))
	router.Handle("/task", middleware.Authenticator(http.HandlerFunc(CreateTask))).Methods("POST")
	router.Handle("/tasks", middleware.Authenticator(http.HandlerFunc(GetTasks)))
	router.Handle("/user", middleware.Authenticator(http.HandlerFunc(DeleteOwnAccount))).Methods("DELETE")
	router.Handle("/users/update", middleware.Authenticator(http.HandlerFunc(UpdateAccount))).Methods("PUT")
	router.Handle("/logout", middleware.Authenticator(http.HandlerFunc(Logout)))

	log.Fatal(http.ListenAndServe(":8080", router))

}