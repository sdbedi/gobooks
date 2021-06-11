package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redeam/gobooks/handlers"
	"github.com/redeam/gobooks/store"
)

// Args args used to run the server
type Args struct {
	// postgres connection string, of the form,
	// e.g "postgres://user:password@localhost:5432/database?sslmode=disable"
	conn string
	// port for the server of the form,
	// e.g ":8080"
	port string
}

// Run run the server based on given args
func Run(args Args) error {
	// router
	router := mux.NewRouter().
		PathPrefix("/api/v1/"). // add prefix for v1 api `/api/v1/`
		Subrouter()

	st := store.NewPostgresBookStore(args.conn)
	hnd := handlers.NewBookHandler(st)
	RegisterAllRoutes(router, hnd)

	// start server
	log.Println("Starting server at port: ", args.port)
	return http.ListenAndServe(args.port, router)
}

// RegisterAllRoutes registers all routes of the api
func RegisterAllRoutes(router *mux.Router, hnd handlers.IBookHandler) {

	// set content type
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	// get books
	router.HandleFunc("/books", hnd.Get).Methods(http.MethodGet)
	// create books
	router.HandleFunc("/books", hnd.Create).Methods(http.MethodPost)
	// delete book
	router.HandleFunc("/books", hnd.Delete).Methods(http.MethodDelete)
	// update book details
	router.HandleFunc("/books/update", hnd.UpdateDetails).Methods(http.MethodPut)
	// list books
	router.HandleFunc("/books/list", hnd.List).Methods(http.MethodGet)
}
