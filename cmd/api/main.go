package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/salilkoirala46/library-management/internal/handlers"
)

type Message struct {
	chats   []string
	friends []string
}

func main() {
	// mux := http.NewServeMux()

	router := mux.NewRouter() //for gorilla mux

	handlers.Handler(router)
	fmt.Println("Starting server on :8000")
	http.ListenAndServe(":8000", router)
}
