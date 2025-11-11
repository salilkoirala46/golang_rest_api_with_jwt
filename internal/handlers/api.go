package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/salilkoirala46/library-management/internal/middleware"
	"github.com/salilkoirala46/library-management/internal/tools"
)

var jwtkey = []byte("secret_key")

func Handler(router *mux.Router) {
	// 👇 Create a subrouter for /User routes
	userRouter := router.PathPrefix("/User").Subrouter()
	userRouter.Use(middleware.AuthorizationMiddleware)
	userRouter.HandleFunc("", GetUsers).Methods("GET")
	userRouter.HandleFunc("", AddUser).Methods("POST")
	userRouter.HandleFunc("/{id}", UpdateUser).Methods("PUT")
	userRouter.HandleFunc("/{id}", DeleteUser).Methods("DELETE")

	router.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "this is the homepage")
	})

	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/register", Register).Methods("POST")

	router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user tools.UserDetail
	var err error
	user = tools.UserDetail{
		Name:  r.FormValue("username"),
		Email: r.FormValue("email"),
	}

	password := r.FormValue("password")

	var database *tools.DatabaseInterface
	database, err = tools.NewDatabase()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var token *tools.Token
	token, err = (*database).RegisterUser(&user, password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user tools.UserDetail
	var err error
	user = tools.UserDetail{
		Name: r.FormValue("username"),
	}

	password := r.FormValue("password")

	var database *tools.DatabaseInterface
	database, err = tools.NewDatabase()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var token *tools.Token
	token, err = (*database).Login(&user, password)
	fmt.Println(err)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if token == nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var err error
	var database *tools.DatabaseInterface
	database, err = tools.NewDatabase()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var usersData *[]tools.UserDetail //variable should start with Capital letter for export
	usersData, err = (*database).GetAllUsers()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usersData)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	var user tools.UserDetail

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var database *tools.DatabaseInterface
	database, err = tools.NewDatabase()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var users *[]tools.UserDetail //variable should start with Capital letter for export

	users, err = (*database).AddUser(&user)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)

}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var updatedUser tools.UserDetail
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var database *tools.DatabaseInterface
	database, err = tools.NewDatabase()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var users *[]tools.UserDetail //variable should start with Capital letter for export

	users, err = (*database).UpdateUser(updatedUser, id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var database *tools.DatabaseInterface
	database, err = tools.NewDatabase()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var users *[]tools.UserDetail //variable should start with Capital letter for export

	users, err = (*database).DeleteUser(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)

}
